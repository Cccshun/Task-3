#include "influence.h"
#include "attack.h"
#include <iostream>




double cal_Influ_Net(int graph[][kN], std::vector<int> seeds) {
	int EdgeNum = 0;
	std::vector<std::vector<int>> g_vec_temp(kN, std::vector<int>());
	for (int i = 0; i < kN; i++)
		for (int j = 0; j < kN; j++) {
			EdgeNum += graph[i][j];
			if (graph[i][j]) g_vec_temp[i].push_back(j);
		}
	EdgeNum /= 2;
	int Att_Sum = (int)(EdgeNum * kAttackPercentTarget);

	if (Att_Sum == 0)
	{
		return cal_Influ(seeds, g_vec_temp);
	}
	else
	{
		auto failureModel = AttackEdge(graph, g_vec_temp, FindMaxLoad, kAttackPercentTarget);
		int num = failureModel.size();
		double fitness = 0;
		for (auto& elem: failureModel) {
			fitness += cal_Influ(seeds, elem);
		}
		return fitness / num;
	}
}


double cal_Influ(std::vector<int> seeds, std::vector<std::vector<int>>& vectorGraph) {
    double influ = 0;
	for (const auto& seed: seeds) {
		influ += 1;
		double sum_temp = 0;
		for (const auto& seedAdjacency: vectorGraph[seed]) {
			sum_temp += 1;
			sum_temp += vectorGraph[seedAdjacency].size() * kActvationProbability;
		}
		influ += kActvationProbability * sum_temp;
	}

	double sec_sum_temp = 0;
	double thr_sum_temp = 0;
	for (int i = 0; i<seeds.size(); i++)
	{
		int seed = seeds[i];
		std::vector<int> Cs = vectorGraph[seed];
		std::vector<int> Cs_simi;
		find_simi(Cs_simi, Cs, seeds);
		double sum_temp = 0;
		for (int j = 0; j<Cs_simi.size(); j++)
		{
			double temp = 0;
			int node = Cs_simi[j];
			temp += 1;
			temp += kActvationProbability * vectorGraph[node].size();
			//temp -= _p;
			//加减都做完之后再乘以系数
			sum_temp += temp * kActvationProbability;
		}
		sec_sum_temp += sum_temp;
		//第二项计算完毕
		std::vector<int> Cs_dis_simi;
		std::vector<int> Cs_d;
		find_third(Cs_dis_simi, Cs_d, Cs, seeds, seed);
		for (int i = 0; i<Cs_dis_simi.size(); i++)
		{
			for (int j = 0; j<Cs_d.size(); j++)
			{
				if (std::find(vectorGraph[Cs_dis_simi[i]].begin(), vectorGraph[Cs_dis_simi[i]].end(), Cs_d[j]) != vectorGraph[Cs_dis_simi[i]].end())
					thr_sum_temp += kActvationProbability * kActvationProbability;
			}
		}
		//第三项计算完毕
	}
	return (influ - sec_sum_temp - thr_sum_temp);
}


void find_simi(std::vector<int>& Cs_simi, std::vector<int> Cs, std::vector<int> seeds)
{
	for (int i = 0; i<seeds.size(); i++)
	{
		int temp_node = seeds[i];
		for (int j = 0; j<Cs.size(); j++)
		{
			if (Cs[j] == temp_node)
				Cs_simi.push_back(temp_node);
		}
	}
}


void find_third(std::vector<int>& Cs_dis_simi, std::vector<int>& Cs_d, std::vector<int> Cs, std::vector<int> seeds, int seed)
{
	std::vector<int> Cs_dis_simi_temp = Cs;
	for (int i = 0; i<Cs_dis_simi_temp.size(); i++)
	{
		for (int j = 0; j<seeds.size(); j++)
			if (Cs_dis_simi_temp[i] == seeds[j])
				Cs_dis_simi_temp[i] = -1;
		if (Cs_dis_simi_temp[i] != -1)
			Cs_dis_simi.push_back(Cs_dis_simi_temp[i]);
	}
	std::vector<int>Cs_simi;
	for (int i = 0; i<seeds.size(); i++)
	{
		int temp_node = seeds[i];
		for (int j = 0; j<Cs.size(); j++)
		{
			if (Cs[j] == temp_node)
				Cs_simi.push_back(temp_node);
		}
	}
	for (int i = 0; i<Cs_simi.size(); i++)
	{
		if (Cs_simi[i] != seed)
			Cs_d.push_back(Cs_simi[i]);
	}
}


double cal_Influ_Model(std::vector<int> seeds, std::vector<std::vector<std::vector<int>>> &failureModel)
{
    int num = failureModel.size();
	double fitness = 0;
	for (auto &elem : failureModel)
	{
		fitness += cal_Influ(seeds, elem);
	}
	return fitness / num;
}
