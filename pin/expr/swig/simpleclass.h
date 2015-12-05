#ifndef SIMPLECLASS_H
#define SIMPLECLASS_H

#include <iostream>
#include <vector>

class SimpleClass
{
public:
    SimpleClass(){};
    std::string hello();
    void helloString(std::vector<std::string> *results);
    void helloBytes(std::vector<char> *results);
    virtual void GetWindow() = 0;
};

class SimpleClassB : SimpleClass
{
    void GetWindow();
};

#endif
