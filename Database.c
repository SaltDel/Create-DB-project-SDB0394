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

/* rebuild memory from log */
void db_init() {
    FILE *file = fopen("data.db", "r");
    if (!file) return;

    char command[10];
    char key[KEY_SIZE];
    char value[VALUE_SIZE];

    while (fscanf(file, "%s %s %s", command, key, value) == 3) {
        if (strcmp(command, "SET") == 0) {
            db_set(key, value);
        }
    }

    fclose(file);
}

/* set key value */
void db_set(const char *key, const char *value) {
    int index = find_key(key);

    if (index >= 0) {
        strcpy(records[index].value, value);
    } else {
        if (record_count < MAX_RECORDS) {
            strcpy(records[record_count].key, key);
            strcpy(records[record_count].value, value);
            record_count++;
        }
    }

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
