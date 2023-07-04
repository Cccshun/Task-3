#include "attack.h"
#include <algorithm>
#include <iterator>
#include <iostream>
#include <fstream>
#include <cmath>


std::vector<std::vector<std::vector<int>>> AttackEdge(int graph[][kN], std::vector<std::vector<int>> vector_graph, std::pair<Edge, Load> (*fp) (const std::map<Edge, Load>&), double per) {
    int sum_edge = 0;
    std::map<Edge, Load> edge_load;  // edgeLoad{<front, back>: <load, capacity>}
    for (int i = 0; i < kN; i++) {
        for (int j = 0; j < i; j++) {
            sum_edge += graph[i][j];
        }

        int degree = vector_graph[i].size();  // degree of node i
        for (const auto j: vector_graph[i]) {
            double load = pow(degree * vector_graph[j].size(), kAlpha);
            edge_load.insert(make_pair(Edge(i, j), Load(load, kBeta * load)));
        }
    }

    int sum_attacked_edge = (int) (per * sum_edge);
    int sum_overload_edge = 0;  // The number of edge that have been overload
    std::vector<std::vector<std::vector<int>>> maliciousattack_model;
    maliciousattack_model.push_back(vector_graph);
    while (sum_overload_edge <= sum_attacked_edge) {
        auto edge = fp(edge_load);
        if (!edge.second.second > 0)
            break;
        DeleteEdge(edge, vector_graph, edge_load);
        LoadRedistribution(edge, vector_graph, edge_load);
        sum_overload_edge++;

        // Cascading failures process
        while (true) {
            auto overload_edge = FindOverloadEdge(edge_load);
            if (overload_edge.second.first == 0 && overload_edge.second.second == 0) 
                break;
            DeleteEdge(overload_edge, vector_graph, edge_load);
            LoadRedistribution(overload_edge, vector_graph, edge_load);
            sum_overload_edge++;
        }
        maliciousattack_model.push_back(vector_graph);
    }
    return maliciousattack_model;
}


void AttackEdgeTest(int graph[][kN], std::pair<Edge, Load> (*fp) (const std::map<Edge, Load>&)) {
    int temp_graph[kN][kN];
    for (int i = 0; i < kN; i++)
        for (int j = 0; j < kN; j++)
            temp_graph[i][j] = graph[i][j];

    // Find the adjacency of nodes and the number of edges
    int sum_edge = 0;
    std::vector<std::vector<int>> adjacency(kN, std::vector<int>());
    for (int i = 0; i < kN; i++)
        for (int j = 0; j < kN; j++)
            if (graph[i][j]) {
                adjacency[i].push_back(j);
                sum_edge += 1;
            }
    sum_edge /= 2;

    std::map<Edge, Load> edge_load;  // edgeLoad{<front, back>: <load, capacity>}
    for (int i = 0; i < kN; i++) {
        int degree = adjacency[i].size();  // degree of node i
        for (const auto j: adjacency[i]) {
            // double load = alpha * degree * adjacency[j].size();
            double load = pow(degree * adjacency[j].size(), kAlpha);
            edge_load.insert(make_pair(Edge(i, j), Load(load, kBeta * load)));
        }
    }

    int sum_attack = (int) (kAttackPercentTarget*sum_edge);
    int attack_count = 0;  // The number of launching attack
    int sum_overload_edge = 0;  // The number of edge that have been overload
    std::ofstream file("../reporting/testData.txt");
    while (sum_overload_edge < sum_attack) {
        auto edge = fp(edge_load);
        if (!edge.second.second > 0) {
            printf("There was no attacking target!");
            break;
        }
        DeleteEdge(edge, adjacency, edge_load);
        LoadRedistribution(edge, adjacency, edge_load);
        attack_count++;
        sum_overload_edge++;
        // printf("\n\n*********************\nAttack edge (%d, %d) load:%f, capacity:%f in the %d-th Step!\n",
            // edge.first.first, edge.first.second, edge.second.first, edge.second.second, attackCount);

        while (true) {
            
            auto overload_edge = FindOverloadEdge(edge_load);
            if (overload_edge.second.first == 0 && overload_edge.second.second == 0) {
                // printf("The overload propagation in the %d-th attack is over!\n\n", attackCount);
                break;
            }
            else {
                // printf("(%d, %d) load: %f, capacity: %f\n",
                //     overloadEdge.first.first, overloadEdge.first.second, overloadEdge.second.first, overloadEdge.second.second);
            }
            DeleteEdge(overload_edge, adjacency, edge_load);
            LoadRedistribution(overload_edge, adjacency, edge_load);
            sum_overload_edge++;
        }
        printf("%f\n", 1 - (double)sum_overload_edge / sum_edge);
        file << 1 - (double) (sum_overload_edge / sum_edge) << std::endl;
    }
    file.close();
    printf("Total number of attack: %d,Total number of overload edge: %d", attack_count, sum_overload_edge);
}


// Find the edge with max load
std::pair<Edge, Load> FindMaxLoad(const std::map<Edge, Load>& el) {
    std::pair<Edge, Load> max = make_pair(Edge(0, 0), Load(0, 0));  // default value
    if (el.size()) {
        for (auto& e: el) {
            if (e.second.first > max.second.first) {
                max = e;
            }
        }
        return max;
    }
    else {
        printf("FindMaxLoad: edge_load was empty!");
        return max;
    }
}


// Random find the edge
std::pair<Edge, Load> RandomFind(const std::map<Edge, Load>& el) {
    std::pair<Edge, Load> edge = make_pair(Edge(0, 0), Load(0, 0));  // default value
    if (el.size()) {
        int index = rand() % el.size();
        auto it = el.begin();
        for (int i = 0; i < index; i++) {
            it++;
        }
        return *it;
    }
    else {
        printf("RandomFind: edge_load was empty!");
        return edge;
    }
}


// Delete the edge (front, back) and (back, front) in edgeLoad
void DeleteEdge(std::pair<Edge, Load> edge, std::vector<std::vector<int>>& adjacency, std::map<Edge, Load>& edge_load) {
    int front = edge.first.first;
    int back = edge.first.second;
    edge_load.erase(std::make_pair(front, back));
    edge_load.erase(std::make_pair(back, front));
    // Delete the neighbor of node front and back
    auto front_to_back = std::find(adjacency[front].begin(), adjacency[front].end(), back);
    auto back_to_front = std::find(adjacency[back].begin(), adjacency[back].end(), front);
    adjacency[front].erase(front_to_back);
    adjacency[back].erase(back_to_front);
}


// Δload(m) = Load(front, back)) * Capactiy(front, m) / (ΣCapacity(front, i) + ΣCapacity(back, j)),   i ∈ Neighbor(front),j ∈ Neighbor(back) 
void LoadRedistribution(std::pair<Edge, Load> edge, std::vector<std::vector<int>>& adjacency, std::map<Edge, Load>& edge_load) {
    int front = edge.first.first;
    int back = edge.first.second;
    double load = edge.second.first;

    // Calculate the sum capacity around the edge (front, back)
    double sum_capacity = 0;
    for (const auto& node: adjacency[front]) {
        sum_capacity += edge_load[std::make_pair(front, node)].second;
    }
    for (const auto& node: adjacency[back]) {
        sum_capacity += edge_load[std::make_pair(back, node)].second;
    }

    // Redistribution the load of edge (front, back) to adjcant edges
    for (int i = 0; i < adjacency[front].size(); i++) {
        int front_neighbor = adjacency[front][i];
        double front_neighbor_capacity = edge_load[std::make_pair(front, front_neighbor)].second;
        edge_load[std::make_pair(front, front_neighbor)].first += load * front_neighbor_capacity / sum_capacity /2;
        edge_load[std::make_pair(front_neighbor, front)].first += load * front_neighbor_capacity / sum_capacity /2;
    }
    for (int i = 0; i < adjacency[back].size(); i++) {
        int backAdj = adjacency[back][i];
        double back_neighbor_capacity = edge_load[std::make_pair(back, backAdj)].second;
        edge_load[std::make_pair(back, backAdj)].first += load * back_neighbor_capacity / sum_capacity /2;
        edge_load[std::make_pair(backAdj, back)].first += load * back_neighbor_capacity / sum_capacity /2;
    }
}


// int CalculateAdjacencySumDegree(int node, std::vector<std::vector<int>>& adjacency) {
//     int SumDegree = 0;
//     for (const auto& item: adjacency[node]) {
//         SumDegree += adjacency[item].size();
//     }
//     return SumDegree;
// }


std::pair<Edge, Load> FindOverloadEdge(const std::map<Edge, Load>& edge_load) {
    for (auto& item: edge_load) {
        if (item.second.first >= item.second.second) return item;
    }

    return std::make_pair(std::make_pair(0, 0), std::make_pair(0, 0)); // No found overload node
}