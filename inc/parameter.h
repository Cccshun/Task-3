#ifndef _PARAMETER_H_
#define _PARAMETER_H_

#include <vector>
#include <algorithm>


//parameters for graphs
const int kN = 1000; // This must be a constant variable or throw mistake
// const int VEC = 0; //Whether using vector or not
const double kAttackPercentTarget = 0.5; // The attack percentage
const double kAttackPercentTest = 0.5;
// const int _Edge = 1;//1-degree, 2-random

//parameters for influence
const int kSeedSize = 20;//Size of seed set
const double kActvationProbability = 0.01;//Activation probability

//parameters for Algorithm
// const int _Target = 2; //1-PlainInflu, 2-AtkNetInflu
// const int initial_Size = 30;
const int kPopSize = 20;
const double kCrossoverProbability = 0.5;
const double kMutationProbability = 0.1;
const double kGlobalSearchProbability = 0.1;
const double kLocalSearchProbability = 0.1;
const int kMaxGenerations = 150;

//parameter for cascading failure
const double kAlpha = 0.5;  // from degree to load of edge
const double kBeta = 1.7;  // from load of edge to capacity of edge

#endif
