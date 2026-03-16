#include <stdio.h>
#include <string.h>
#include "database.h"

int main() {

    char command[10];
    char key[100];
    char value[100];

    db_init();

    while (scanf("%9s", command) != EOF) {

        if (strcmp(command, "SET") == 0) {

            if (scanf("%99s %99s", key, value) == 2) {
                db_set(key, value);
            }

        }
        else if (strcmp(command, "GET") == 0) {

            if (scanf("%99s", key) == 1) {

                char *result = db_get(key);

                if (result)
                    printf("%s\n", result);
                else
                    printf("NULL\n");

                fflush(stdout);
            }

        }
        else if (strcmp(command, "EXIT") == 0) {

            break;
        }
    }

    return 0;
}
