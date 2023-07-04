#ifndef _ATTACK_H_
#define _ATTACK_H_

#include "parameter.h"
#include <map>


 
typedef std::pair<int, int> Edge;  // <front, back>
typedef std::pair<double, double> Load;  // <load, capacity>

std::vector<std::vector<std::vector<int>>> AttackEdge(int graph[][kN], std::vector<std::vector<int>> vectorGraph, std::pair<Edge, Load> (*fp) (const std::map<Edge, Load>&), double);
void AttackEdgeTest(int graph[][kN], std::pair<Edge, Load> (*fp) (const std::map<Edge, Load>&));

std::pair<Edge, Load> FindMaxLoad(const std::map<Edge, Load>&);
std::pair<Edge, Load> RandomFind(const std::map<Edge, Load>&);

void LoadRedistribution(std::pair<Edge, Load> edge, std::vector<std::vector<int>>& adjacency, std::map<Edge, Load>& edgeLoad);
void DeleteEdge(std::pair<Edge, Load> edge, std::vector<std::vector<int>>& adjacency, std::map<Edge, Load>& edgeLoad);

std::pair<Edge, Load> FindOverloadEdge(const std::map<Edge, Load>& edgeLoad);

int CalculateAdjacencySumDegree(int node, std::vector<std::vector<int>>& adjacency);
#endif