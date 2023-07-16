#include <bits/stdc++.h>

using namespace std;

int main(){
    ios::sync_with_stdio(false);
    cin.tie(0);

    int q;
    cin >> q;

    while (q--){
        int l, r;
        cin >> l >> r;

        int64_t pr = 1;
        for (int i = l; i <= r; i++){
            pr *= i;
        }

        int64_t sum;
        do {
            sum = 0;
            while (pr > 0){
                sum += pr%10;
                pr /= 10;
            }
            pr = sum;
        } while (sum > 9);
        cout << sum << '\n';
    }
}