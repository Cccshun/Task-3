#ifndef _ATTACK_H_
#define _ATTACK_H_

#include "parameter.h"
#include <map>

using namespace std;

 
typedef pair<int, int> Edge;  // <front, back>
typedef pair<double, double> Load;  // <load, capacity>

vector<vector<vector<int>>> AttackEdge(int graph[][kN], vector<vector<int>> vectorGraph, pair<Edge, Load> (*fp) (const map<Edge, Load>&), double);
void AttackEdgeTest(int graph[][kN], pair<Edge, Load> (*fp) (const map<Edge, Load>&));

pair<Edge, Load> FindMaxLoad(const map<Edge, Load>&);
pair<Edge, Load> RandomFind(const map<Edge, Load>&);

void LoadRedistribution(pair<Edge, Load> edge, vector<vector<int>>& adjacency, map<Edge, Load>& edgeLoad);
void DeleteEdge(pair<Edge, Load> edge, vector<vector<int>>& adjacency, map<Edge, Load>& edgeLoad);

pair<Edge, Load> FindOverloadEdge(const map<Edge, Load>& edgeLoad);

int CalculateAdjacencySumDegree(int node, vector<vector<int>>& adjacency);
#endif