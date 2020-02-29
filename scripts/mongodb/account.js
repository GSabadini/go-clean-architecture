db = db.getSiblingDB('bank');

db.createUser({
    user: 'dev',
    pwd: 'dev',
    roles: [
        {
            role: 'root',
            db: 'admin',
        },
    ],
});

accounts = db.createCollection('accounts');
db.accounts.createIndex( { "cpf": 1 }, { unique: true } )

db.createCollection('transfers');
