#include "genetic.h"

#include <cstdlib>
#include <iostream>
#include <cmath>
#include <algorithm>
#include <iterator>
#include <fstream>


std::ostream_iterator<int> outIter(std::cout, " ");


void GA::Init() {
    for (int i = 0; i < kPopSize; i++) {
        for (int j = 0; j < kSeedSize; j++) {
            m_pop[i].X.push_back(_produce_random());
        }
        _remove_duplivate(m_pop[i].X);
    }
    _load_vector_graph();
    m_maliciousAttackModel = AttackEdge(G, m_vectorGraph, FindMaxLoad, kAttackPercentTarget);
    m_maliciousAttackModelTest = AttackEdge(G, m_vectorGraph, FindMaxLoad, kAttackPercentTest);

}


void GA::Crossover() {
    std::random_shuffle(m_pop.begin(), m_pop.end());
    m_popChild.resize(kPopSize);  // resize enough size
    for (int i = 0; i < kPopSize; i+=2) {
        for (int j = 0; j < kSeedSize; j++) {
            if ( (rand() / (double)RAND_MAX) < kCrossoverProbability) {
                m_popChild[i].X.push_back(m_pop[i+1].X[j]);
                m_popChild[i+1].X.push_back(m_pop[i].X[j]);
            }
            else {
                m_popChild[i].X.push_back(m_pop[i].X[j]);
                m_popChild[i+1].X.push_back(m_pop[i+1].X[j]);
            }
        }
        _remove_duplivate(m_popChild[i].X);
        _remove_duplivate(m_popChild[i+1].X);
    }
}


void GA::Mutation() {
    for (int i = 0; i < kPopSize; i++) {
        if ( (rand() / (double)RAND_MAX) < kMutationProbability) {
            int tempSeed = _produce_random();
            std::vector<int>::const_iterator cb = m_popChild[i].X.cbegin();
            std::vector<int>::const_iterator ce = m_popChild[i].X.cend();
            while (find(cb, ce, tempSeed) != ce) {
                tempSeed = _produce_random();
            }

            int pos = rand() % kSeedSize;
            m_popChild[i].X[pos] = tempSeed;
        }
    }
}


void GA::SearchLocal() {
    for (auto& ind: m_popChild) {
        if ( rand() / (double)RAND_MAX ) {
            for (int i = 0; i < ind.X.size(); i++) {
                int bestGene = ind.X[i];
                double bestFitness = ind.fitness;
                for (const auto gene: m_vectorGraph[ind.X[i]]) {
                    if (std::find(ind.X.begin(), ind.X.end(), gene) == ind.X.end()) {
                        ind.X[i] = gene;
                        // _Evaluate(G, ind);
                        _MaliciousModelEvaluate(ind);
                        if (ind.fitness > bestFitness) {
                            bestGene = gene;
                            bestFitness = ind.fitness;
                        }
                    }
                }
                ind.X[i] = bestGene;
                ind.fitness = bestFitness;
            }
        }
    }
}


void GA::Select() {
    m_popChild.insert(m_popChild.end(), m_pop.begin(), m_pop.end()); // Merge pop and popchild to popchild
    std::sort(m_popChild.begin(), m_popChild.end(), [](Seeds ind1, Seeds ind2) { return ind1.fitness > ind2.fitness; });
    std::copy(m_popChild.begin(), m_popChild.begin() + kPopSize, m_pop.begin());  // Best seeds save to pop
    m_popChild.clear();
}


void GA::Evaluation() {
    for (int i = 0; i < kPopSize; i++) {
    // _Evaluate(G, m_pop[i]);
    // _Evaluate(G, m_popChild[i]);
    _MaliciousModelEvaluate(m_pop[i]);
    _MaliciousModelEvaluate(m_popChild[i]);
    }
}


void GA::FindBest(std::string fileName) {
    Init();
    int gen = 0;
    std::ofstream file;
    file.open(fileName);
    // fileMalicious.open("C:\\Users\\Lenovo\\Desktop\\ResearchTask\\Task2\\reporting\\Malicious.txt");
    // std::ofstream fileRandom;
    // fileRandom.open("C:\\Users\\Lenovo\\Desktop\\ResearchTask\\Task2\\reporting\\Random.txt");
    while (gen++ < kMaxGenerations) {
        Crossover();
        Mutation();
        Evaluation();
        Select();
        printf("\nGen-%d  fit:%f  ", gen, m_pop[0].fitness);
        // fileMalicious << m_pop[0].fitness << std::endl;
        std::copy(m_pop[0].X.begin(), m_pop[0].X.end(), outIter);
        printf("  Test:%f", cal_Influ_Model(m_pop[0].X, m_maliciousAttackModelTest));
        // fileRandom << _random_attack_evaluate(m_pop[0]) << std::endl;
        file << m_pop[0].fitness << std::endl;
    }
    file.close();
    // fileMalicious.close();
    // fileRandom.close();
}


void GA::_remove_duplivate(std::vector<int>& seeds) {
    std::sort(seeds.begin(), seeds.end());
    seeds.erase(std::unique(seeds.begin(), seeds.end()), seeds.end());

    for (int i = seeds.size(); i < kSeedSize; i++) {
        auto seed = _produce_random();
        if (std::find(seeds.begin(), seeds.end(), seed) == seeds.end()) {
            seeds.push_back(seed);
        }
    }
}


void GA::_load_vector_graph() {
    for (int i = 0; i < kN; i++) {
            for (int j = 0; j < kN; j++) {
                if (G[i][j])
                    m_vectorGraph[i].push_back(j);
        }
    }
}


double GA::_random_attack_evaluate(Seeds& seeds) {
    auto randomAttackModel = AttackEdge(G, m_vectorGraph, RandomFind, kAttackPercentTarget);
    return cal_Influ_Model(seeds.X, randomAttackModel);
}


void MA::FindBest(std::string fileName) {
    Init();
    int gen = 0;
    std::ofstream file;
    file.open(fileName);
    // fileMalicious.open("C:\\Users\\Lenovo\\Desktop\\ResearchTask\\Task2\\reporting\\Malicious.txt");
    // std::ofstream fileRandom;
    // fileRandom.open("C:\\Users\\Lenovo\\Desktop\\ResearchTask\\Task2\\reporting\\Random.txt");
    while (gen++ < kMaxGenerations) {
        Crossover();
        Mutation();
        Evaluation();
        SearchLocal();
        Select();
        printf("\nGen-%d  fit:%f  ", gen, m_pop[0].fitness);
        // fileMalicious << m_pop[0].fitness << std::endl;
        std::copy(m_pop[0].X.begin(), m_pop[0].X.end(), outIter);
        printf("  Test:%f", cal_Influ_Model(m_pop[0].X, m_maliciousAttackModelTest));
        // fileRandom << _random_attack_evaluate(m_pop[0]) << std::endl;
        file << m_pop[0].fitness << std::endl;
    }
    file.close();
    // fileMalicious.close();
    // fileRandom.close();

}