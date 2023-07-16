#include <bits/stdc++.h>

using namespace std;

typedef struct {
    int urods;
    int power;
} stable;

typedef struct {
    int a;
    int b;
    int c;
    int power;
} tabun;

int max3(int a, int b, int c){
    return max(max(a, b), c);
}

int min3(int a, int b, int c){
    return min(min(a, b), c);
}

string bin(uint n) {
    string bin = "";
    while(n > 0) {
        int nn = n % 2;
        bin += (char)(nn + '0'); 
        n /= 2;
    }
    return bin;
}

int main(){
    ios::sync_with_stdio(false);
    cin.tie(0);

    int n;
    cin >> n;

    vector<tabun> tabuns(n); 
    string strtmp;
    tabun tabtmp;
    for (int i = 0; i < n; i++){
        tabtmp = tabun{0, 0, 0, 0};
        cin >> strtmp;
        for (int j = 0; j < strtmp.length(); j++){
            if (strtmp[j] == 'a') tabtmp.a++;
            if (strtmp[j] == 'b') tabtmp.b++;
            if (strtmp[j] == 'c') tabtmp.c++;
        }
        tabtmp.power = strtmp.length();
        tabuns[i] = tabtmp;
    }

    stable minUrods = {numeric_limits<int>::max(), 0};
    stable kombo = {0, 0};

    for (uint64_t i = 1; i < pow(2, n); i++){
        vector<tabun> curStable(0);
        string mask = bin(i);
        for (int j = 0; j < mask.length(); j++){
            if (mask[j] == '1'){
                curStable.push_back(tabuns[j]);
            }
        }
        int a = 0, b = 0, c = 0, power = 0;
        for (int j = 0; j < curStable.size(); j++){
            a+=curStable[j].a;
            b+=curStable[j].b;
            c+=curStable[j].c;
            power+=curStable[j].power;
        }
        // cout << mask << ":  " << a << ' ' << b << ' ' << c << ' ' << power << endl;
        kombo.power = power;
        kombo.urods = max3(a, b, c) - min3(a, b, c);
        if (kombo.urods < minUrods.urods) minUrods = kombo;
        else if (kombo.urods == minUrods.urods){
            if (kombo.power > minUrods.power) minUrods.power = kombo.power;
        } 
    }
    cout << minUrods.power << '\n';
}