#ifndef DATABASE_H
#define DATABASE_H

#define MAX_RECORDS 1000
#define KEY_SIZE 100
#define VALUE_SIZE 100

typedef struct {
    char key[KEY_SIZE];
    char value[VALUE_SIZE];
} Record;

void db_init();
void db_set(const char *key, const char *value);
char* db_get(const char *key);

#endif
