#include <stdio.h>
#include <string.h>
#include "database.h"

static Record records[MAX_RECORDS];
static int record_count = 0;

 //search for key index 
static int find_key(const char *key) {
    for (int i = 0; i < record_count; i++) {
        if (strcmp(records[i].key, key) == 0) {
            return i;
        }
    }
    return -1;
}

/* internal set used when rebuilding log */
static void db_set_memory(const char *key, const char *value) {

    int index = find_key(key);

    if (index >= 0) {

        strncpy(records[index].value, value, VALUE_SIZE - 1);
        records[index].value[VALUE_SIZE - 1] = '\0';

    } else {

        if (record_count < MAX_RECORDS) {

            strncpy(records[record_count].key, key, KEY_SIZE - 1);
            records[record_count].key[KEY_SIZE - 1] = '\0';

            strncpy(records[record_count].value, value, VALUE_SIZE - 1);
            records[record_count].value[VALUE_SIZE - 1] = '\0';

            record_count++;
        }
    }
}

/* rebuild memory from log */
void db_init() {

    FILE *file = fopen("data.db", "a+");   // ensures file exists
    if (!file) return;

    rewind(file);

    char command[10];
    char key[KEY_SIZE];
    char value[VALUE_SIZE];

    while (fscanf(file, "%9s %99s %99s", command, key, value) == 3) {

        if (strcmp(command, "SET") == 0) {
            db_set_memory(key, value);
        }
    }

    fclose(file);
}

/* set key value */
void db_set(const char *key, const char *value) {

    db_set_memory(key, value);

    /* append to disk */
    FILE *file = fopen("data.db", "a");

    if (file) {
        fprintf(file, "SET %s %s\n", key, value);
        fclose(file);
    }
}

/* get value */
char* db_get(const char *key) {

    int index = find_key(key);

    if (index >= 0) {
        return records[index].value;
    }

    return NULL;
}
