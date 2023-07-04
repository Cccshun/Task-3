#ifndef _SA_H
#define _SA_H
#include "genetic.h"


class SA: GA
{
private:
    double m_t;
    double m_delta = 0.1;
    int m_gen = 0;
    std::vector<int> m_seed;
public:
    SA(double t): m_t(t) { };
    void find(std::string);
private:
    void init();
};

#endif