db = db.getSiblingDB('stone');

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

db.createCollection('stone');
