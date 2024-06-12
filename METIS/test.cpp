#include <metis.h>
#include <vector>
#include <iostream>
#include <fstream>
#include <string>
#include <sstream>
using namespace std;

vector<idx_t> func(vector<idx_t> xadj, vector<idx_t> adjncy, vector<idx_t> vwgt) {
	idx_t nVertices = xadj.size() - 1; // 节点数
	idx_t nEdges = adjncy.size() / 2;// 边数
	idx_t nWeights = 1;
	idx_t nParts = 2;// 子图个数
	idx_t objval;
	std::vector<idx_t> part(nVertices, 0);
        int ret = METIS_PartGraphKway(&nVertices, &nWeights, xadj.data(), adjncy.data(),
		vwgt.data(), NULL, NULL, &nParts, NULL,
		NULL, NULL, &objval, part.data());
	std::cout << ret << std::endl;
	for (unsigned part_i = 0; part_i < part.size(); part_i++) {
		std::cout << part_i << " " << part[part_i] << std::endl;
	}
	return part;
}

int main() {
	string infile = "graph.txt";
	string outfile = "partition.txt";
	ifstream ingraph(infile);
	if (!ingraph) {
		cout << "打开文件失败！" << endl;
		exit(1);//失败退回操作系统   
	}
	cout << "打开文件 successed！" << endl;
	int vexnum, edgenum;
	string line;
	getline(ingraph, line);
	istringstream tmp(line);
	tmp >> vexnum >> edgenum;
	vector<idx_t> xadj(0);
	vector<idx_t> adjncy(0); //点的id从0开始
	vector<idx_t> vwgt(0);
	idx_t a, w;
	for (int i = 0; i < vexnum; i++) {
		xadj.push_back(adjncy.size());
		getline(ingraph, line);
		istringstream tmp(line);
		while (tmp >> a >> w) {
			adjncy.push_back(a);
			vwgt.push_back(w);
		}
	}
	xadj.push_back(adjncy.size());
	ingraph.close();
	vector<idx_t> part = func(xadj, adjncy, vwgt);
	ofstream outpartition(outfile);
	if (!outpartition) {
		cout << "打开文件失败！" << endl;
		exit(1);
	}
	for (int i = 0; i < part.size(); i++) {
		outpartition << i << " " << part[i] << endl;
	}
	outpartition.close();
	return 0;
}
