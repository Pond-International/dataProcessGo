package repositories

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/zap"
	"pondDataProcessGo/models"
	"pondDataProcessGo/utils"
	"strings"
)

const (
	url        = "bolt://memgraph-1198751972.us-east-1.elb.amazonaws.com:7687"
	dbUsername = "pond"
	dbPassword = "DiveIntoPond"
)

func connectGraphDb() (*neo4j.DriverWithContext, context.Context, error) {
	driver, err := neo4j.NewDriverWithContext(url, neo4j.BasicAuth(dbUsername, dbPassword, ""))
	ctx := context.Background()
	if err != nil {
		zap.L().Error("connectGraghDbError", zap.Errors("err", []error{err}))
	}
	return &driver, ctx, err
}

type GraphRepository struct {
	graphDb *neo4j.DriverWithContext
	ctx     context.Context
}

func NewGraphRepository() *GraphRepository {
	db, ctx, err := connectGraphDb()
	if err != nil {
		zap.L().Error("NewGraphRepository", zap.Errors("err", []error{err}))
	}
	return &GraphRepository{
		graphDb: db,
		ctx:     ctx,
	}
}

func (r *GraphRepository) MergeTwitterAccount(users []models.User) error {
	for i := 0; i < len(users); i++ {
		user := users[i]
		cypher := `
			MERGE (n:TwitterAccount {username: "%s"})
			SET n.followerCount = %d,
				n.followingCount = %d,
				n.id = "%s",
				n.name = "%s"
    `
		username := strings.ToLower(user.Username)
		followerCount := user.PublicMetrics.FollowersCount
		followingCount := user.PublicMetrics.FollowingCount
		id := user.ID
		name := user.Name
		query := fmt.Sprintf(cypher, username, followerCount, followingCount, id, name)
		zap.L().Info("MergeTwitterAccount_cypher", zap.String("cypher", query))
		_, err := neo4j.ExecuteQuery(r.ctx, *r.graphDb, query, nil, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
		if err != nil {
			zap.L().Error("MergeTwitterAccount", zap.Errors("err", []error{err}))
			return err
		}
	}
	return nil
}

func (r *GraphRepository) MergeTwitterFollowing(user models.User, followUsers []models.User, uType bool) error {
	//followerUsers可以是follower也可以是followering
	username := strings.ToLower(user.Username)

	for i := 0; i < len(followUsers); i++ {
		cypher := fmt.Sprintf(`
		MATCH (source:TwitterAccount {username: "%s"})
		MATCH (target:TwitterAccount {username: "%s"})
		MERGE (source) -[e:TwitterFollowing]-> (target)
	`, strings.ToLower(followUsers[i].Username), username)
		if uType == false {
			//很关键,follower和following 顺序是相反的
			cypher = fmt.Sprintf(`
		MATCH (source:TwitterAccount {username: "%s"})
		MATCH (target:TwitterAccount {username: "%s"})
		MERGE (source) -[e:TwitterFollowing]-> (target)
	`, username, strings.ToLower(followUsers[i].Username))
		}
		//add log
		if i == 0 {
			zap.L().Info("MergeTwitterFollowing_cypher", zap.String("cypher", cypher))
		}
		_, err := neo4j.ExecuteQuery(r.ctx, *r.graphDb, cypher, nil, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
		if err != nil {
			zap.L().Error("MergeTwitterFollowing", zap.Errors("err", []error{err}))
			return err
		}
	}
	return nil
}

func (r *GraphRepository) MergeTwitter2Person(newPersons []models.User) []int64 {
	var nodesToAdd []int64
	for i := 0; i < len(newPersons); i++ {
		cypher := fmt.Sprintf(`
		MATCH (n:TwitterAccount {username: "%s"})
		MERGE (n)-[:Twitter2Person]->(p:Person)
		RETURN id(p)
	`, strings.ToLower(newPersons[i].Username))
		if i == 0 {
			zap.L().Info("MergeTwitter2Person_cypher", zap.String("cypher", cypher))
		}
		result, err := neo4j.ExecuteQuery(r.ctx, *r.graphDb, cypher, nil, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
		if err != nil {
			zap.L().Error("MergeTwitter2Person", zap.Errors("err", []error{err}))
		}
		priId := result.Records[0].AsMap()["id(p)"].(int64)
		nodesToAdd = append(nodesToAdd, priId)
	}
	return nodesToAdd
}

func (r *GraphRepository) MergePerson2Person(user models.User, rType bool) {
	//如果是r1.following 则传true f2.following传false
	useType := "r1"
	if rType == false {
		useType = "r2"
	}
	cypher := fmt.Sprintf(`MATCH (:TwitterAccount {username: "%s"})-[:Twitter2Person]->(p1:Person)
MATCH (:TwitterAccount {username: "%s"})<-[:TwitterFollowing]-(:TwitterAccount)-[:Twitter2Person]->(p2:Person)
MERGE (p1)-[r1:Person2Person]->(p2)
MERGE (p2)-[r2:Person2Person]->(p1) 
SET %s.following=1`, strings.ToLower(user.Username), strings.ToLower(user.Username), useType)
	zap.L().Info("MergePerson2Person_cypher", zap.String("cypher", cypher))

	_, err := neo4j.ExecuteQuery(r.ctx, *r.graphDb, cypher, nil, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		zap.L().Error("MergePerson2Person", zap.Errors("err", []error{err}))
	}
}

func (r *GraphRepository) UpdateComputeRelationStrength(user models.User) ([]int64, []int64, []float64) {
	cypher := fmt.Sprintf(`
	MATCH (t1:TwitterAccount {username: "%s"})-[:Twitter2Person]->(p1:Person)
	MATCH (p1)-[r1:Person2Person]->(p2:Person)
	MATCH (p2)-[r2:Person2Person]->(p1)
	MATCH (t2:TwitterAccount)-[:Twitter2Person]->(p2:Person)
	WITH t1, t2, p1, p2, r1, r2
	OPTIONAL MATCH (l1:LensAccount)-[:Lens2Person]->(p1)
	OPTIONAL MATCH (l2:LensAccount)-[:Lens2Person]->(p2)
	WITH 1.0 as w, 1.0 as lw, r1, r2, p1, p2,
		log(coalesce(t1.followerCount, 0) + 1) as t1FeC,
		log(coalesce(t2.followerCount, 0) + 1) as t2FeC,
		log(coalesce(t1.followingCount, 0) + 1) as t1FiC,
		log(coalesce(t2.followingCount, 0) + 1 ) as t2FiC,
		log(coalesce(l1.followerCount, 0) + 1) as l1FeC,
		log(coalesce(l2.followerCount,0) + 1) as l2FeC,
		log(coalesce(l1.followingCount, 0) + 1) as l1FiC,
		log(coalesce(l2.followingCount, 0) + 1) as l2FiC
	WITH CASE WHEN t1FiC > 0 AND t2FeC > 0 THEN (1 + sqrt(t1FeC)) * (1 + atan(w/t1FiC/t2FeC)) ELSE 0 END AS A1,
		CASE WHEN t2FiC > 0 AND t1FeC > 0 and r2.following = 1 then (1 + sqrt(t2FeC)) * (1 + atan(w/t2FiC/t1FeC)) else 0 end as A2,
case when l1FiC > 0 and l2FeC > 0 and r1.lensFollowing = 1 then (1 + sqrt(l1FeC)) * ( 1 + atan(lw/l1FiC/l2FeC)) else 0 end as B1,
	case when l2FiC > 0 and l1FeC > 0 and r2.lensFollowing = 1 then (1 + sqrt(l2FeC)) * (1 + atan(lw/l2FiC/l1FeC)) else 0 end as B2,
	atan(coalesce(r1.transferTo, 0)) as C1,
	atan(coalesce(r2.transferTo, 0)) as C2,
	r1, r2, p1, p2
	WITH exp(- 0.05 * A1 - 0.05 * A2 - 0.05 * B1 - 0.05 * B2 - 0.3 * C1 - 0.3 * C2) as distance, r1, r2, p1, p2
	SET r1.distance = distance, r2.distance = distance
	RETURN collect(id(p1)) as sources, collect(id(p2)) as targets, collect(r1.distance) as distances
	`, strings.ToLower(user.Username))
	zap.L().Info("UpdateComputeRelationStrength_cypher", zap.String("cypher", cypher))

	result, err := neo4j.ExecuteQuery(r.ctx, *r.graphDb, cypher, nil, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		zap.L().Error("UpdateComputeRelationStrength", zap.Errors("err", []error{err}))
	}
	sources := result.Records[0].AsMap()["sources"].([]interface{})
	targets := result.Records[0].AsMap()["targets"].([]interface{})
	distances := result.Records[0].AsMap()["distances"].([]interface{})
	var sourcesSlice []int64
	var targetsSlice []int64
	var distanceSlice []float64
	for i := 0; i < len(sources); i++ {
		sourcesSlice = append(sourcesSlice, sources[i].(int64))
	}
	for i := 0; i < len(targets); i++ {
		targetsSlice = append(targetsSlice, targets[i].(int64))
	}
	for i := 0; i < len(distances); i++ {
		distanceSlice = append(distanceSlice, distances[i].(float64))
	}
	return sourcesSlice, targetsSlice, distanceSlice
}

func (r *GraphRepository) AddNodesFromWithNodes(nodes []int64) {
	cypher := `
	CALL kssp.add_nodes($nodesToAdd)
	YIELD nodes_added
	RETURN nodes_added
	`
	zap.L().Info("AddNodesFromWithNodes_cypher", zap.String("cypher", cypher), zap.String("params", utils.Int64SliceToStringLimit5(nodes)), zap.Int("paramsLengths", len(nodes)))
	_, err := neo4j.ExecuteQuery(r.ctx, *r.graphDb, cypher, map[string]interface{}{"nodesToAdd": nodes}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		zap.L().Error("AddNodesFromWithNodes", zap.Errors("err", []error{err}))
	}
}

func (r *GraphRepository) AddNodesFromSourcesTargetsWeights(sources []int64, targets []int64, weights []float64) {
	cypher := `
	CALL kssp.add_edges($sources, $targets, $weights)
	YIELD edges_added
	RETURN edges_added
	`
	zap.L().Info("AddNodesFromSourcesTargetsWeights_cypher", zap.String("cypher", cypher), zap.String("sourcesParams", utils.Int64SliceToStringLimit5(sources)), zap.Int("sourcesParamsLengths", len(sources)), zap.String("targetsParams", utils.Int64SliceToStringLimit5(targets)), zap.Int("targetsParamsLengths", len(targets)), zap.String("weightParams", utils.Float64SliceToStringLimit5(weights)), zap.Int("weightsParamsLengths", len(weights)))

	_, err := neo4j.ExecuteQuery(r.ctx, *r.graphDb, cypher, map[string]interface{}{
		"sources": sources,
		"targets": targets,
		"weights": weights,
	}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		zap.L().Error("AddNodesFromSourcesTargetsWeights", zap.Errors("err", []error{err}))
	}

}

func (r *GraphRepository) GetTwitterAccountInfo(name string) int64 {
	cypher := fmt.Sprintf(`MATCH (p:TwitterAccount {username: "%s"}) return id(p)`, strings.ToLower(name))
	zap.L().Info("GetTwitterAccountInfo_cypher", zap.String("cypher", cypher))
	ret, err := neo4j.ExecuteQuery(r.ctx, *r.graphDb, cypher, nil, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		zap.L().Error("GetTwitterAccountInfo", zap.Errors("err", []error{err}))
	}
	return ret.Records[0].AsMap()["id(p)"].(int64)
	//fmt.Println(ret.Records[0].AsMap()[""].(int64))
}

func (r *GraphRepository) EasyTest() {
	user := models.User{}
	user.Username = "colinpond5362"
	r.MergePerson2Person(user, true)
	//session := (*r.graphDb).NewSession(r.ctx, neo4j.SessionConfig{DatabaseName: ""})
	//defer session.Close(r.ctx)
	//cypher := "MATCH (p:Person) return id(p),p limit 1"
	//cypher := "MATCH (p:Person) WHERE (id(p)>= 3684988 AND id(p)<=3684999) return collect(id(p)) as sources"
	//result, err := neo4j.ExecuteQuery(r.ctx, *r.graphDb, cypher, nil, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	//if err != nil {
	//	panic(err)
	//}
	//sources := result.Records[0].AsMap()["sources"].([]interface{})
	//for i := 0; i < len(sources); i++ {
	//	fmt.Println(sources[i].(int64))
	//}
	//fmt.Println(result.Records[0].AsMap()["collect(id(p))"].([]int))
	//result, err := session.Run(r.ctx, cypher, nil)
	//fmt.Println(result.Record())
	//fmt.Println(err)
	//if err != nil {
	//	log.Fatal("Query execution failed: ", err)
	//}
	//for result.Next(r.ctx) {
	//	record := result.Record()
	//	fmt.Println(record.Get("sources"))
	//}
}
