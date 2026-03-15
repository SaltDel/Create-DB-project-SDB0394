#include <stdio.h>
#include <string.h>
#include "database.h"

int main() {
    char line[256];

    db_init();

    while (fgets(line, sizeof(line), stdin)) {

        char command[10];
        char key[100];
        char value[100];

        if (sscanf(line, "%s", command) != 1)
            continue;

        if (strcmp(command, "SET") == 0) {

            if (sscanf(line, "%s %s %s", command, key, value) == 3) {
                db_set(key, value);
            }

        } else if (strcmp(command, "GET") == 0) {

            if (sscanf(line, "%s %s", command, key) == 2) {
                char *result = db_get(key);

                if (result)
                    printf("%s\n", result);
                else
                    printf("NULL\n");

                fflush(stdout);
            }

        } else if (strcmp(command, "EXIT") == 0) {

            break;

        }
    }

    return 0;
}
