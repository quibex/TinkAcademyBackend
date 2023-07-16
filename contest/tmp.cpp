#include <iostream>
#include <string>
#include <vector>
#include <math.h>

using namespace std;

struct herd {
    int a;
    int b;
    int c;
};

void horses_count(herd &h, string &s) {
    for (int i = 0; i < s.size(); i++)
    {
        if(s[i] == 'a') {
            h.a++;
        }
        else if(s[i] == 'b') {
            h.b++;
        }
        else if(s[i] == 'c') {
            h.c++;
        }
    }
}

string to_bin(int n) {
    string bin = "";
    while(n > 0) {
        int rem = n % 2;
        bin = (char)(rem + '0') + bin; 
        n /= 2;
    }
    return bin;
}

int main() {

    int n;
    cin >> n;
    vector<herd> ar;
    for(int i = 0; i < n; i++) {
        string s;
        cin >> s;
        herd cur_herd = {0, 0, 0};
        horses_count(cur_herd, s);
        ar.push_back(cur_herd);
    }
    // for (int i = 0; i < n; i++) {
    //     cout << ar[i].a << " " << ar[i].b << " " << ar[i].c << endl;    
    // }
    

    int min_ugly = 1e9;
    int max_power = -1e9;
    for (int i = 1; i < pow(2, n); i++) {
        string mask = to_bin(i);
        while(mask.size() != n) {
            mask = "0" + mask;
        }
        // cout << mask << endl;
        herd h = {0, 0, 0};
        for (int j = 0; j < mask.size(); j++) {
            if(mask[j] == '1') {
                h = {h.a + ar[j].a, h.b + ar[j].b, h.c + ar[j].c};
                // cout << h.a << " " << h.b << " " << h.c << endl; 
                int cur_ugly = max(max(h.a, h.b), h.c) - min(min(h.a, h.b), h.c);
                int cur_power = h.a + h.b + h.c;
                if(cur_ugly < min_ugly) {
                    min_ugly = cur_ugly;
                    if(cur_power > max_power) {
                        max_power = cur_power;
                    }
                }
            }
        }
    }
    cout << max_power << endl;

    return 0;
}