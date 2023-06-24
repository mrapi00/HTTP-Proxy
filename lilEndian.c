#include <stdio.h>

int main() {
    int x = 1; // ['0x1', '0x0', '0x0', '0x0'
    char *LSB = (char*) &x; // either ['0x1'] or ['0x0']
    char c = *LSB;
    printf("%c", c);
}