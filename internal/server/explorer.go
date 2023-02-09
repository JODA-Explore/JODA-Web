package server

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/JODA-Explore/BETZE/generator"
	"github.com/JODA-Explore/BETZE/languages"
	extjoda "github.com/JODA-Explore/BETZE/languages/joda"
	"github.com/JODA-Explore/BETZE/query"
)

func explorerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		explorerHandleForm(w, r)
		return
	}
	data := createData("BETZE Benchmark Generator", "flask")
	templateFile := "explorer.html"

	data["timestamp"] = time.Now().UnixNano()

	data["predicates"] = generator.GetPredicateFactoryRepo().GetAllIDs()
	data["aggregations"] = generator.GetAggregationFactoryRepo().GetAllIDs()

	data["languages"] = languages.LanguageIndex()

	sources, err := JodaInstance.GetSources()
	if err != nil {
		addError(data, err)
		executeTemplate(w, templateFile, data)
		return
	}
	datasets := []string{}
	for _, source := range sources {
		datasets = append(datasets, source.Name)
	}
	data["datasets"] = datasets

	executeTemplate(w, templateFile, data)
}

func explorerHandleForm(w http.ResponseWriter, r *http.Request) {
	data := createData("Explorer Benchmark Generator", "compass")
	errorTemplateFile := "explorer.html"
	templateFile := "explorerresults.html"

	config := generatorFromForm(r, data)
	if config == nil {
		executeTemplate(w, errorTemplateFile, data)
		return
	}

	queries, err := generateQuerySet(*config)
	if err != nil {
		addError(data, err)
		executeTemplate(w, errorTemplateFile, data)
		return
	}
	// Statistics
	data["statistics"] = config.QueryGenerator.Statistics()
	data["config"] = config.QueryGenerator.PrintConfig()

	//Network
	network := config.QueryGenerator.Network()
	convertedNetwork := generatorNetworkToD3(network)
	data["graph"] = convertedNetwork

	link_query_map := make(map[int]int)
	query_index := 0
	for i, l := range convertedNetwork.Edges {
		if l.JumpType != 3 {
			continue
		}
		link_query_map[i] = query_index
		query_index++
	}

	data["linkQuery"] = link_query_map

	// Queries
	language_map := make(map[string][]string)

	// Betze Format
	b, err := generator.MarshalQueries(queries, config.QueryGenerator.PrintConfig())
	if err != nil {
		addError(data, err)
	} else {
		language_map["Betze"] = []string{string(b)}
	}

	var chosen_languages []languages.Language
	for _, lang := range config.Languages {
		// Get language from languages.LanguageIndex() where lang equals language.ShortName()
		for _, l := range languages.LanguageIndex() {
			if l.ShortName() == lang {
				chosen_languages = append(chosen_languages, l)
				break
			}
		}
	}

	var intermediate_language []languages.Language
	var non_intermediate_language []languages.Language
	for _, lang := range chosen_languages {
		if config.Intermediate && lang.SupportsIntermediate() {
			intermediate_language = append(intermediate_language, lang)
		} else {
			non_intermediate_language = append(non_intermediate_language, lang)
		}
	}

	for _, l := range intermediate_language {
		language_map[l.Name()] = []string{}
		// Translate queries
		for _, q := range queries {
			language_map[l.Name()] = append(language_map[l.Name()], l.Translate(q))
		}
	}

	queries = query.RemoveIntermediateSets(queries)

	for _, l := range non_intermediate_language {
		language_map[l.Name()] = []string{}
		// Translate queries
		for _, q := range queries {
			language_map[l.Name()] = append(language_map[l.Name()], l.Translate(q))
		}
	}

	data["queries"] = language_map
	executeTemplate(w, templateFile, data)
}

type formConfig struct {
	QueryGenerator *generator.Generator
	Datasets       []string
	Intermediate   bool
	CheckQueries   bool
	Languages      []string
	NumQueries     int
}

func generatorFromForm(r *http.Request, data map[string]interface{}) *formConfig {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		addError(data, err)
		return nil
	}

	valid := true
	//Datasets
	datasets := r.Form["dataset"]
	//Seed
	seed_str := r.Form.Get("seed")
	seed, err := strconv.ParseInt(seed_str, 10, 64)
	if err != nil {
		addError(data, err)
		valid = false
	}
	//# Queries
	num_queries_str := r.Form.Get("num_queries")
	num_queries, err := strconv.Atoi(num_queries_str)
	if err != nil {
		addError(data, err)
		valid = false
	}
	//Min Selectivity
	min_selectivity_str := r.Form.Get("min_selectivity")
	min_selectivity, err := strconv.ParseFloat(min_selectivity_str, 64)
	if err != nil {
		addError(data, err)
		valid = false
	}

	//Max Selectivity
	max_selectivity_str := r.Form.Get("max_selectivity")
	max_selectivity, err := strconv.ParseFloat(max_selectivity_str, 64)
	if err != nil {
		addError(data, err)
		valid = false
	}

	//Probability Backtrack
	probability_backtrack_str := r.Form.Get("probability_backtrack")
	probability_backtrack, err := strconv.ParseFloat(probability_backtrack_str, 64)
	if err != nil {
		addError(data, err)
		valid = false
	}

	//Random jump
	probability_randomjump_str := r.Form.Get("probability_randomjump")
	probability_randomjump, err := strconv.ParseFloat(probability_randomjump_str, 64)
	if err != nil {
		addError(data, err)
		valid = false
	}

	//Intermediate
	intermediate_str := r.Form.Get("intermediate_sets")
	intermediate := intermediate_str == "on"

	//Aggregation
	aggregate_str := r.Form.Get("aggregate")
	aggregate := aggregate_str == "on"
	if aggregate {
		intermediate = false
	}

	//Check Queries
	check_queries_str := r.Form.Get("check_queries")
	check_queries := check_queries_str == "on"

	//Languages
	languages := r.Form["language"]

	// Predicates
	factory := generator.GetPredicateFactoryRepo()
	predicates := r.Form["predicates"]
	if len(predicates) == 0 {
		factory.SetDefault()
	} else {
		for _, pred := range predicates {
			err := factory.Include(pred)
			if err != nil {
				addWarnMessage(data, err)
			}
		}
	}

	// Aggregations
	agg_factory := generator.GetAggregationFactoryRepo()
	aggregations := r.Form["aggregations"]
	if len(aggregations) == 0 {
		agg_factory.SetDefault()
	} else {
		for _, pred := range aggregations {
			err := agg_factory.Include(pred)
			if err != nil {
				addWarnMessage(data, err)
			}
		}
	}

	if !valid {
		return nil
	}

	query_generator := generator.New(seed)
	query_generator.MinSelectivity = min_selectivity
	query_generator.MaxSelectivity = max_selectivity
	query_generator.RandomBrowseProb = probability_randomjump
	query_generator.GoBackProb = probability_backtrack
	query_generator.Predicates = factory.GetChosen()
	if aggregate {
		query_generator.Aggregations = agg_factory.GetChosen()
		query_generator.AggregationProb = 1.0
	}

	return &formConfig{
		QueryGenerator: &query_generator,
		Datasets:       datasets,
		Intermediate:   intermediate,
		CheckQueries:   check_queries,
		Languages:      languages,
		NumQueries:     num_queries,
	}
}

func generateQuerySet(config formConfig) ([]query.Query, error) {
	// Connect to JODA
	con, err := extjoda.Connect(JodaInstance.GetHost())
	if err != nil {
		return nil, err
	}

	// Get datasets
	datasets, err := con.GetDatasets(config.Datasets)
	if err != nil {
		return nil, err
	}

	if config.CheckQueries {
		queries, err := config.QueryGenerator.GenerateQuerySetWithJoda(datasets, int64(config.NumQueries), *con)
		if err != nil {
			return nil, err
		}
		return queries, nil
	}

	queries := config.QueryGenerator.GenerateQuerySet(datasets, int64(config.NumQueries))
	return queries, nil
}

func generatorNetworkToD3(n generator.Network) explorerD3Network {
	edges := []explorerD3Edge{}
	for _, e := range n.Edges {
		if e.From == "" || e.To == "" || e.From == e.To || e.JumpType == 0 {
			continue
		}
		edge := explorerD3Edge{
			From:      e.From,
			To:        e.To,
			JumpType:  e.JumpType,
			Timestamp: e.Timestamp,
		}
		if !e.Query.IsCopy() {
			edge.QueryString = e.Query.String()
		}
		edges = append(edges, edge)
	}

	nodes := []explorerD3Node{}
	for _, e := range n.Nodes {
		group := 0
		if e.Original {
			group = 1
		}
		nodes = append(nodes, explorerD3Node{
			Name:      e.DSName,
			Group:     group,
			Size:      e.Size,
			Timestamp: e.Timestamp,
		})
	}
	return explorerD3Network{
		Nodes: nodes,
		Edges: edges,
		MaxTS: n.MaxTimestamp,
	}
}

type explorerD3Node struct {
	Name      string `json:"id"`
	Group     int    `json:"group"`
	Size      uint64 `json:"size"`
	Timestamp uint   `json:"timestamp"`
}
type explorerD3Edge struct {
	From        string `json:"source"`
	To          string `json:"target"`
	QueryString string `json:"query"`
	JumpType    int    `json:"jump"`
	Timestamp   uint   `json:"timestamp"`
}
type explorerD3Network struct {
	Nodes []explorerD3Node `json:"nodes"`
	Edges []explorerD3Edge `json:"links"`
	MaxTS uint             `json:"max_ts"`
}
