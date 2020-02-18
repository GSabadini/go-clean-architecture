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

db.createCollection('accounts');
db.createCollection('transfers');
