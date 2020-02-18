db = db.getSiblingDB('account');

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

db.createCollection('account');
