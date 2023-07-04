#ifndef _INFLUENCE_H_
#define _INFLUENCE_H_

#include "parameter.h"

double cal_Influ_Net(int graph[][kN], std::vector<int> seeds);

double cal_Influ(std::vector<int> seeds, std::vector<std::vector<int>>& vectorGraph);

void find_simi(std::vector<int>& Cs_simi, std::vector<int> Cs, std::vector<int> seeds);
void find_third(std::vector<int>& Cs_dis_simi, std::vector<int>& Cs_d, std::vector<int> Cs, std::vector<int> seeds, int seed);

double cal_Influ_Model(std::vector<int> seeds, std::vector<std::vector<std::vector<int>>>&);

#endif