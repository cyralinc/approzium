#include <iostream>
#include "dbauth_libpq.h"

using namespace std;

int main() {
    char *s = dbauth_get_hashed_password("SALT");
    cout << "sheesh" << endl;    
}
