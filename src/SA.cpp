#include "SA.h"
#include <algorithm>
#include <cmath>
#include <fstream>


void SA::init() {
    for (int i = 0; i < kSeedSize; i++)
    {
        int s = ProduceRandom();
        m_seed.push_back(s);
        for(int j = 0; j < i; j++) {
            if (m_seed[j] == m_seed[i])
            {
                 i--;
                 m_seed.pop_back();
                 break;
            }     
        }
    }

    LoadVectorGraph();
    malicious_attack_model_ = AttackEdge(G, vector_graph_, FindMaxLoad, kAttackPercentTarget);
    malicious_attack_model_test_ = AttackEdge(G, vector_graph_, FindMaxLoad, kAttackPercentTest); 
}


void SA::find(std::string fileName) {
    init();
    int count = 0;
    std::ofstream file;
    file.open(fileName);
    while (m_t > 0.000001)
    {
        count++;
       // generate a random gene;
       int newGene = ProduceRandom();
       for (int i = 0; i < kSeedSize; i++)
       {
            if (newGene == m_seed[i])
            {
                i = 0;
                newGene = ProduceRandom();
            }
       }
       
        // generate a exchange position;
        int index = rand() % kSeedSize;
        int oldGene = m_seed[index];
        double oldFitness = cal_Influ_Model(m_seed, malicious_attack_model_);

        // replace old gene with new gene if newFitness > oldFitness
        m_seed[index] = newGene;
        double newFitness = cal_Influ_Model(m_seed, malicious_attack_model_);
        if (oldFitness > newFitness)
            m_seed[index] = oldGene;

        if (((double) rand() / RAND_MAX) < std::exp(- m_delta / m_t))
            m_seed[index] = newGene;
        

        m_t *= 0.91;
        
        printf("count:%d  t:%f  probability:%f, fitness:%f, newfitness:%f  ", 
            count, m_t, std::exp(- m_delta / m_t), cal_Influ_Model(m_seed, malicious_attack_model_), newFitness);
        
        for (int i = 0; i < kSeedSize; i++)
            std::cout << m_seed[i] << " ";
        std::cout << std::endl;
        file << cal_Influ_Model(m_seed, malicious_attack_model_) << std::endl;
    }
    file.close();
}