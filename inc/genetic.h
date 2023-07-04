#ifndef _GENETIC_H_
#define _GENETIC_H_

#include <iostream>
#include <memory>
#include "parameter.h"
#include "influence.h"
#include "attack.h"

extern int G[][kN];

struct Seeds
{
    std::vector<int> X;
    double fitness;
};


class GA {
    friend double cal_Influ_Net(int graph[][kN], std::vector<int>& seeds);
    friend std::vector<std::vector<std::vector<int>>> AttackEdge(int graph[][kN], std::vector<std::vector<int>>& vectorGraph);
    friend std::pair<Edge, Load> FindMaxLoad(const std::map<Edge, Load>&);
    friend std::pair<Edge, Load> RandomFind(const std::map<Edge, Load>&);
    // friend double cal_Influ(std::vector<int>& seeds, std::vector<std::vector<int>>& g_vec, int graph[][N]);
    // friend void find_simi(std::vector<int>& Cs_simi, std::vector<int> Cs, std::vector<int> seeds);
    // friend void find_third(std::vector<int>& Cs_dis_simi, std::vector<int>& Cs_d, std::vector<int> Cs, std::vector<int> seeds, int seed);

public:
    GA(): m_pop(std::vector<Seeds> (kPopSize)), m_popChild(std::vector<Seeds> (kPopSize)),
        m_vectorGraph(std::vector<std::vector<int>> (kN, std::vector<int>()))
        {
            printf("**: N: %d, MAX_GEN: %d, POP_SIZE: %d, SEEDS_SIZE:%d, pc: %f, pm: %f, pgs: %f, pls: %f\n",
                kN, kMaxGenerations, kPopSize, kSeedSize, kCrossoverProbability, kMutationProbability, kGlobalSearchProbability, kLocalSearchProbability);
        }

public:
    void FindBest(std::string);

protected:
    void Init();
    void Crossover();
    void Mutation();
    void SearchGlobal();
    void SearchLocal();
    void Select();
    void Evaluation();
    // void _Evaluate(int graph[][N], Seeds& seeds) {
    //     // seeds.fitness = cal_Influ_Net(graph, seeds.X);
    //     seeds.fitness = cal_Influ(seeds.X, m_vectorGraph);
    // }
    void _MaliciousModelEvaluate(Seeds& seeds) { seeds.fitness = cal_Influ_Model(seeds.X, m_maliciousAttackModel); };
    double _random_attack_evaluate(Seeds& seeds);

protected:
    void _load_vector_graph();
    int _produce_random() { return rand() % kN; }
    void _remove_duplivate(std::vector<int>&);

protected:
    std::vector<Seeds> m_pop;
    std::vector<Seeds> m_popChild;

protected:
    std::vector<std::vector<int>> m_vectorGraph;
    std::vector<std::vector<std::vector<int>>> m_maliciousAttackModel;
    std::vector<std::vector<std::vector<int>>> m_maliciousAttackModelTest;
};


class MA: GA {
public:
    void FindBest(std::string); 
};

#endif