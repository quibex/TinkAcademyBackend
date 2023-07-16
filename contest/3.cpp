#include <bits/stdc++.h>

using namespace std;

int main(){
    ios::sync_with_stdio(false);
    cin.tie(0);

    int n, m, q;
    cin >> n >> m >> q;

    vector<vector<int>> arr(n, vector<int>(m));

    for (int i = 0; i < n; i++){
        for (int j = 0; j < m; j++){
            cin >> arr[i][j];
        }
    }

    int x, y, k;
    while (q--){
        cin >> x >> y >> k;
        x--; y--;

        int ans = 0;
        for (int i = 0; i < n; i++){
            if (i == x) continue;
            if (abs(arr[i][y] - arr[x][y]) <= k) ans++;
        }
        for (int j = 0; j < m; j++){
            if (j == y) continue;
            if (abs(arr[x][j] - arr[x][y]) <= k) ans++;
        }
        cout << ans << '\n';
    }

}