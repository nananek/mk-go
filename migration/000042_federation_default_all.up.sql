-- Switch the federation column default from 'none' to 'all' so fresh mk-go
-- installs don't end up with federation completely disabled now that #536
-- enforces the field on deliver / inbox paths. Misskey TS の admin UI
-- placeholder / 既定値も 'all' で、ここを揃えると drop-in 互換が取れる。
ALTER TABLE "meta" ALTER COLUMN "federation" SET DEFAULT 'all';

-- Existing rows with the legacy default (= 'none') are residual values
-- carried over from when this column had no enforcement. Flip them to
-- 'all' so existing mk-go databases keep federating after upgrade. 既に
-- admin が "none" を意図して保存している運用は事実上存在しないため
-- (pre-#536 では UI 入力経路があっても挙動に影響しなかった) この一括
-- update は安全。
UPDATE "meta" SET "federation" = 'all' WHERE "federation" = 'none';
