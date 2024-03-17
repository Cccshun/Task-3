#include "genetic.h"

#include <cstdlib>
#include <iostream>
#include <cmath>
#include <algorithm>
#include <iterator>
#include <fstream>


ostream_iterator<int> outIter(cout, " ");


void GA::Init() {
    for (int i = 0; i < kPopSize; i++) {
        for (int j = 0; j < kSeedSize; j++) {
            pop_[i].X.push_back(ProduceRandom());
        }
        RemoveDuplication(pop_[i].X);
    }
    LoadVectorGraph();
    malicious_attack_model_ = AttackEdge(G, vector_graph_, FindMaxLoad, kAttackPercentTarget);
    malicious_attack_model_test_ = AttackEdge(G, vector_graph_, FindMaxLoad, kAttackPercentTest);

}


void GA::Crossover() {
    random_shuffle(pop_.begin(), pop_.end());
    pop_offspring_.resize(kPopSize);  // resize enough size
    for (int i = 0; i < kPopSize; i+=2) {
        for (int j = 0; j < kSeedSize; j++) {
            if ( (rand() / (double)RAND_MAX) < kCrossoverProbability) {
                pop_offspring_[i].X.push_back(pop_[i+1].X[j]);
                pop_offspring_[i+1].X.push_back(pop_[i].X[j]);
            }
            else {
                pop_offspring_[i].X.push_back(pop_[i].X[j]);
                pop_offspring_[i+1].X.push_back(pop_[i+1].X[j]);
            }
        }
        RemoveDuplication(pop_offspring_[i].X);
        RemoveDuplication(pop_offspring_[i+1].X);
    }
}


void GA::Mutate() {
    for (int i = 0; i < kPopSize; i++) {
        if ( (rand() / (double)RAND_MAX) < kMutationProbability) {
            int tempSeed = ProduceRandom();
            vector<int>::const_iterator cb = pop_offspring_[i].X.cbegin();
            vector<int>::const_iterator ce = pop_offspring_[i].X.cend();
            while (find(cb, ce, tempSeed) != ce) {
                tempSeed = ProduceRandom();
            }

            int pos = rand() % kSeedSize;
            pop_offspring_[i].X[pos] = tempSeed;
        }
    }
}


void GA::SearchLocal() {
    for (auto& ind: pop_offspring_) {
        if ( rand() / (double)RAND_MAX ) {
            for (int i = 0; i < ind.X.size(); i++) {
                int bestGene = ind.X[i];
                double bestFitness = ind.fitness;
                for (const auto gene: vector_graph_[ind.X[i]]) {
                    if (std::find(ind.X.begin(), ind.X.end(), gene) == ind.X.end()) {
                        ind.X[i] = gene;
                        // _Evaluate(G, ind);
                        MaliciousModelEvaluate(ind);
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
    pop_offspring_.insert(pop_offspring_.end(), pop_.begin(), pop_.end()); // Merge pop and popchild to popchild
    sort(pop_offspring_.begin(), pop_offspring_.end(), [](Seeds ind1, Seeds ind2) { return ind1.fitness > ind2.fitness; });
    copy(pop_offspring_.begin(), pop_offspring_.begin() + kPopSize, pop_.begin());  // Best seeds save to pop
    pop_offspring_.clear();
}


void GA::Evaluate() {
    for (int i = 0; i < kPopSize; i++) {
    // _Evaluate(G, m_pop[i]);
    // _Evaluate(G, m_popChild[i]);
    MaliciousModelEvaluate(pop_[i]);
    MaliciousModelEvaluate(pop_offspring_[i]);
    }
}


void GA::FindBest(std::string fileName) {
    Init();
    int gen = 0;
    ofstream file;
    file.open(fileName);
    // fileMalicious.open("C:\\Users\\Lenovo\\Desktop\\ResearchTask\\Task2\\reporting\\Malicious.txt");
    // std::ofstream fileRandom;
    // fileRandom.open("C:\\Users\\Lenovo\\Desktop\\ResearchTask\\Task2\\reporting\\Random.txt");
    while (gen++ < kMaxGenerations) {
        Crossover();
        Mutate();
        Evaluate();
        Select();
        printf("\nGen-%d  fit:%f  ", gen, pop_[0].fitness);
        // fileMalicious << m_pop[0].fitness << std::endl;
        copy(pop_[0].X.begin(), pop_[0].X.end(), outIter);
        printf("  Test:%f", cal_Influ_Model(pop_[0].X, malicious_attack_model_test_));
        // fileRandom << _random_attack_evaluate(m_pop[0]) << std::endl;
        file << pop_[0].fitness << std::endl;
    }
    file.close();
    // fileMalicious.close();
    // fileRandom.close();
}


void GA::RemoveDuplication(std::vector<int>& seeds) {
    sort(seeds.begin(), seeds.end());
    seeds.erase(unique(seeds.begin(), seeds.end()), seeds.end());

    for (int i = seeds.size(); i < kSeedSize; i++) {
        auto seed = ProduceRandom();
        if (find(seeds.begin(), seeds.end(), seed) == seeds.end()) {
            seeds.push_back(seed);
        }
    }
}


void GA::LoadVectorGraph() {
    for (int i = 0; i < kN; i++) {
            for (int j = 0; j < kN; j++) {
                if (G[i][j])
                    vector_graph_[i].push_back(j);
        }
    }
}


double GA::RandomAttackEvaluate(Seeds& seeds) {
    auto randomAttackModel = AttackEdge(G, vector_graph_, RandomFind, kAttackPercentTarget);
    return cal_Influ_Model(seeds.X, randomAttackModel);
}


void MA::FindBest(string fileName) {
    Init();
    int gen = 0;
    ofstream file;
    file.open(fileName);
    // fileMalicious.open("C:\\Users\\Lenovo\\Desktop\\ResearchTask\\Task2\\reporting\\Malicious.txt");
    // ofstream fileRandom;
    // fileRandom.open("C:\\Users\\Lenovo\\Desktop\\ResearchTask\\Task2\\reporting\\Random.txt");
    while (gen++ < kMaxGenerations) {
        Crossover();
        Mutate();
        Evaluate();
        SearchLocal();
        Select();
        printf("\nGen-%d  fit:%f  ", gen, pop_[0].fitness);
        // fileMalicious << m_pop[0].fitness << endl;
        copy(pop_[0].X.begin(), pop_[0].X.end(), outIter);
        printf("  Test:%f", cal_Influ_Model(pop_[0].X, malicious_attack_model_test_));
        // fileRandom << _random_attack_evaluate(m_pop[0]) << endl;
        file << pop_[0].fitness << endl;
    }
    file.close();
    // fileMalicious.close();
    // fileRandom.close();

}