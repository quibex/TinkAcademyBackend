#include <bits/stdc++.h>

using namespace std;

int main(){
    ios::sync_with_stdio(false);
    cin.tie(0);

    int n;
    cin >> n;

    vector<int> shurs(n, 0);

    int shurn = 1;
    int64_t shursum = 0;
    int64_t res = 0;

    for (int i = 0; i < n; i++){

        shurn = (2*(i+1)-1);
        shursum += pow(shurn, 2);

        res = pow(shurn, 3) - shursum;
        cout << res << ' ';
    }

}