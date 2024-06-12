#include "metis.h"

#include <vector>

#include <iostream>

#include <fstream>

#include <string>

#include <sstream>

using namespace std;

vector<idx_t> func(vector<idx_t> xadj, vector<idx_t> adjncy, vector<idx_t> vwgt, vector<idx_t> wwgt, idx_t &fmt, int nParts) {

	idx_t nVertices = xadj.size() - 1; // 节点数

	idx_t nEdges = adjncy.size() / 2;// 边数

	idx_t nWeights = 1;

	// idx_t nParts = 8;// 子图个数

	idx_t objval;

	std::vector<idx_t> part(nVertices, 0);


	idx_t options[METIS_NOPTIONS]; 
	METIS_SetDefaultOptions(options); 
	options[METIS_OPTION_NO2HOP] = 1; 
	options[METIS_OPTION_NITER] = 1;
	options[METIS_OPTION_CTYPE] = METIS_CTYPE_RM;
	options[METIS_OPTION_IPTYPE] = METIS_IPTYPE_EDGE;
	options[METIS_OPTION_RTYPE] = METIS_RTYPE_SEP1SIDED;
// options
	int ret = METIS_PartGraphKway(&nVertices, &nWeights, xadj.data(), adjncy.data(),

		wwgt.data(), NULL, vwgt.data(), &nParts, NULL,

		NULL, options, &objval, part.data());
	 
	// int ret = METIS_PartGraphRecursive(&nVertices, &nWeights, xadj.data(), adjncy.data(),

	// 	wwgt.data(), NULL, vwgt.data(), &nParts, NULL,

	// 	NULL, options, &objval, part.data());

	std::cout << ret << std::endl;

	// for (unsigned part_i = 0; part_i < part.size(); part_i++) {

	// 	std::cout << part_i << " " << part[part_i] << std::endl;

	// }

	return part;

}

int main(int argc, char* argv[]) {
	cout<<"yesssssssssssssssssss";

	ifstream ingraph(argv[1]); //C:\\Users\\admin\\Desktop\\graph partition\\test\\Project1\\
	
	if (!ingraph) {

		cout << "打开文件失败！" << endl;

		exit(1);//失败退回操作系统   

	}

	int vexnum, edgenum;

	string line;

	getline(ingraph, line);

	istringstream tmp(line);
	idx_t fmt;

	tmp >> vexnum >> edgenum;

	vector<idx_t> xadj(0);

	vector<idx_t> adjncy(0); //点的id从0开始

	vector<idx_t> vwgt(0);

	vector<idx_t> wwgt(0);

	idx_t a, w, z;

	for (int i = 0; i < vexnum; i++) {

		xadj.push_back(adjncy.size());

		getline(ingraph, line);

		istringstream tmp(line);
		tmp >> z;
		wwgt.push_back(z);
		while (tmp >> a >> w) {

			adjncy.push_back(a);

			vwgt.push_back(w);

		}

	}

	xadj.push_back(adjncy.size());

	ingraph.close();

	//for (int i = 0; i < wwgt.size(); i++) {
	//	cout << i << " " << wwgt[i] << endl;
	//}
	//cout << endl;
	//for (int i = 0; i < vwgt.size(); i++) {
	//	cout << i << " " << vwgt[i] << endl;
	//}
	//cout << endl;
	//for (int i = 0; i < xadj.size(); i++) {
	//	cout << i << " " << xadj[i] << endl;
	//}
	//cout << endl;
	//for (int i = 0; i < adjncy.size(); i++) {
	//	cout << i << " " << adjncy[i] << endl;
	//}
	int nParts = atoi(argv[3]);
	vector<idx_t> part = func(xadj, adjncy, vwgt, wwgt, fmt, nParts);

	ofstream outpartition(argv[2]); //C:\\Users\\admin\\Desktop\\graph partition\\test\\Project1\\

	if (!outpartition) {

		cout << "打开文件失败！" << endl;

		exit(1);

	}

	for (int i = 0; i < part.size(); i++) {

		outpartition << i << " " << part[i] << endl;

	}

	outpartition.close();

}
