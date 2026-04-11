print("Starting migration to internal_user_id...");
print("This migration creates UserAccount for each existing user");
print("and updates their training_sessions to use internal_user_id\n");

print("1. Creating indexes for user_accounts...");
db.user_accounts.createIndex({ "id": 1 }, { unique: true });
db.user_accounts.createIndex(
    { "platform": 1, "platform_user_id": 1 },
    { unique: true }
);
db.user_accounts.createIndex({ "internal_user_id": 1 });
print("✓ Indexes created\n");

print("2. Creating indexes for link_codes...");
db.link_codes.createIndex({ "id": 1 }, { unique: true });
db.link_codes.createIndex({ "code": 1 });
db.link_codes.createIndex(
    { "expires_at": 1 },
    { expireAfterSeconds: 0 }
);
print("✓ Indexes created\n");

print("3. Migrating existing training sessions and logs...");

const uniqueUserIds = db.training_sessions.distinct("user_id");
print(`Found ${uniqueUserIds.length} unique users\n`);

let migratedSessionsCount = 0;
let migratedLogsCount = 0;
let skippedCount = 0;

uniqueUserIds.forEach(function (userId) {
    const isUUID = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(userId);

    if (isUUID) {
        skippedCount++;
        return;
    }

    let platform = "telegram";
    const userIdNum = parseInt(userId);

    if (!isNaN(userIdNum) && userIdNum < 100000000) {
        platform = "vk";
    }

    function generateUUID() {
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
            const r = Math.random() * 16 | 0;
            const v = c === 'x' ? r : (r & 0x3 | 0x8);
            return v.toString(16);
        });
    }

    const internalUserId = generateUUID();
    const accountId = generateUUID();
    const now = new Date();

    const account = {
        id: accountId,
        internal_user_id: internalUserId,
        platform: platform,
        platform_user_id: userId,
        username: "",
        linked_at: now,
        created_at: now
    };

    try {
        db.user_accounts.insertOne(account);

        const sessionsResult = db.training_sessions.updateMany(
            { "user_id": userId },
            { $set: { "user_id": internalUserId } }
        );

        const logsResult = db.training_logs.updateMany(
            { "user_id": userId },
            { $set: { "user_id": internalUserId } }
        );

        migratedSessionsCount += sessionsResult.modifiedCount;
        migratedLogsCount += logsResult.modifiedCount;
        print(`  ✓ ${platform} user ${userId} -> ${internalUserId} (${sessionsResult.modifiedCount} sessions, ${logsResult.modifiedCount} logs)`);
    } catch (e) {
        print(`  ✗ Error for user ${userId}: ${e.message}`);
    }
});

print(`\n✓ Migrated ${migratedSessionsCount} sessions and ${migratedLogsCount} logs`);
print(`  Skipped ${skippedCount} already migrated\n`);

print("4. Verification...");

const totalSessions = db.training_sessions.countDocuments();
const sessionsWithUUID = db.training_sessions.countDocuments({
    "user_id": { $regex: /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i }
});

const totalLogs = db.training_logs.countDocuments();
const logsWithUUID = db.training_logs.countDocuments({
    "user_id": { $regex: /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i }
});

print(`Total sessions: ${totalSessions}`);
print(`Sessions with UUID: ${sessionsWithUUID}`);
print(`Total logs: ${totalLogs}`);
print(`Logs with UUID: ${logsWithUUID}`);

if (sessionsWithUUID === totalSessions && logsWithUUID === totalLogs) {
    print("✓ All sessions and logs migrated!\n");
} else {
    if (sessionsWithUUID !== totalSessions) {
        print(`⚠ Warning: ${totalSessions - sessionsWithUUID} sessions not migrated`);
    }
    if (logsWithUUID !== totalLogs) {
        print(`⚠ Warning: ${totalLogs - logsWithUUID} logs not migrated`);
    }
    print("");
}

print("5. Statistics:");
print(`Total UserAccounts: ${db.user_accounts.countDocuments()}`);
print(`Telegram: ${db.user_accounts.countDocuments({ platform: "telegram" })}`);
print(`VK: ${db.user_accounts.countDocuments({ platform: "vk" })}`);

print("\n✓ Migration completed!");
print("\nIMPORTANT: Existing users will see their data correctly");
print("because UserAccount was created with their platform_user_id");
