#ifndef _GENETIC_H_
#define _GENETIC_H_

#include <iostream>
#include <memory>
#include "parameter.h"
#include "influence.h"
#include "attack.h"

using namespace std;
extern int G[][kN];

struct Seeds
{
    std::vector<int> X;
    double fitness;
};


class GA {
    friend double cal_Influ_Net(int graph[][kN], vector<int>& seeds);
    friend vector<vector<vector<int>>> AttackEdge(int graph[][kN], vector<vector<int>>& vectorGraph);
    friend pair<Edge, Load> FindMaxLoad(const map<Edge, Load>&);
    friend pair<Edge, Load> RandomFind(const map<Edge, Load>&);
    // friend double cal_Influ(std::vector<int>& seeds, std::vector<std::vector<int>>& g_vec, int graph[][N]);
    // friend void find_simi(std::vector<int>& Cs_simi, std::vector<int> Cs, std::vector<int> seeds);
    // friend void find_third(std::vector<int>& Cs_dis_simi, std::vector<int>& Cs_d, std::vector<int> Cs, std::vector<int> seeds, int seed);

public:
    GA(): pop_(vector<Seeds> (kPopSize)), pop_offspring_(vector<Seeds> (kPopSize)),
        vector_graph_(vector<vector<int>> (kN, vector<int>()))
        {
            printf("**: N: %d, MAX_GEN: %d, POP_SIZE: %d, SEEDS_SIZE:%d, pc: %f, pm: %f, pgs: %f, pls: %f\n",
                kN, kMaxGenerations, kPopSize, kSeedSize, kCrossoverProbability, kMutationProbability, kGlobalSearchProbability, kLocalSearchProbability);
        }

public:
    void FindBest(string);

protected:
    void Init();
    void Crossover();
    void Mutate();
    void SearchGlobal();
    void SearchLocal();
    void Select();
    void Evaluate();
    // void _Evaluate(int graph[][N], Seeds& seeds) {
    //     // seeds.fitness = cal_Influ_Net(graph, seeds.X);
    //     seeds.fitness = cal_Influ(seeds.X, m_vectorGraph);
    // }
    void MaliciousModelEvaluate(Seeds& seeds) { seeds.fitness = cal_Influ_Model(seeds.X, malicious_attack_model_); };
    double RandomAttackEvaluate(Seeds& seeds);

protected:
    void LoadVectorGraph();
    int ProduceRandom() { return rand() % kN; }
    void RemoveDuplication(vector<int>&);

protected:
    vector<Seeds> pop_;
    vector<Seeds> pop_offspring_;

protected:
    vector<vector<int>> vector_graph_;
    vector<vector<vector<int>>> malicious_attack_model_;
    vector<vector<vector<int>>> malicious_attack_model_test_;
};


class MA: GA {
public:
    void FindBest(string); 
};

#endif