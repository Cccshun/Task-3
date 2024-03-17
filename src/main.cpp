#include <iostream>
#include <string>
#include <stdlib.h>
#include <ctime>
#include <cstring>
#include <fstream>
#include <set>
#include "genetic.h"
#include "SA.h"
#include "influence.h"
#include "parameter.h"
#include "attack.h"
#include <unistd.h>

using namespace std;

int G[kN][kN];
void LoadGraph(int graph[][kN], string filepath);

int main()
{
  srand(time(0)); 
  string filepath = "../network/BA1000B.txt";
  // string filepath = "../network/WS1000B.txt";
  // string filepath = "../network/ER1000B.txt";
  LoadGraph(G, filepath);

  string network = "BA1000";
  // string network = "WS1000";
  // string network = "ER1000";

  string filename_1 = "../data1000/"+ network +"-ga"+ to_string(kSeedSize) +".txt";
  class GA ga;
  ga.FindBest(filename_1); 

  string filename_2 = "../data1000/"+ network +"-ma"+ to_string(kSeedSize) +".txt";
  class MA ma;
  ma.FindBest(filename_2); 

  string filename_3 = "../data1000/"+ network +"-sa"+ to_string(kSeedSize) +".txt";
  class SA sa(1);
  sa.find(filename_3);
}


void LoadGraph(int graph[][kN], string filepath) 
{
  ifstream fp;
  fp.open(filepath);
  if (!fp) {
    printf("Fail to open the file!");
    return;
  } else {
    for (int i = 0; i < kN; ++i)
    {
      char temp[kN+1];
      fp.getline(temp, kN+1);
      for (int j = 0; j < kN; j++)
      {
        G[i][j] = (int)temp[j] - 48;
      }
    }
    fp.close();
  }
}