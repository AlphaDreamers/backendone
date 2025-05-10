* **Number of SQL tables defined:** 18

    1. Badge
    2. UserBadge
    3. User
    4. Skill
    5. GigTag
    6. user\_skills
    7. Biometrics
    8. Gig
    9. RegistrationToken
    10. GigImage
    11. GigPackage
    12. GigPackageFeature
    13. Review
    14. Category
    15. Order
    16. Chat
    17. Message
    18. \_GigToGigTag

* **Number of Go model types (excluding the `OrderStatus` enum):** 16

---

## Missing or extra models

| SQL table             | Go struct model present? | Notes                                                                        |
| --------------------- | ------------------------ | ---------------------------------------------------------------------------- |
| **Badge**             | ❌ missing                | No `Badge` struct.                                                           |
| **UserBadge**         | ❌ missing                | No `UserBadge` struct.                                                       |
| **RegistrationToken** | ❌ missing                | No model for registration tokens.                                            |
| **Country**           | 🔶 extra                 | Go adds a `Country` model but there is no `Country` table in the SQL schema. |

---

## Field‑level and relation mismatches

1. **User ↔ Country**

    * SQL: `User` has a `country` TEXT column.
    * Go: `User` refers to `CountryID uuid.UUID` and a `Country` struct.

2. **User**

    * SQL: `firstName` + `lastName`;
    * Go: single `Name` + `Surname`.
    * SQL: `twoFactorVerified` vs Go: `KYCVerified`.

3. **Skill / user\_skills**

    * SQL: `Skill.id` is TEXT (UUID), `user_skills.skillId` TEXT.
    * Go: `SkillId int` on `UserSkill` (type mismatch) and references `SkillID` (note uppercase “ID”) not `skillId`.

4. **Biometrics**

    * SQL: table named `Biometrics`, columns `id`, `type`, `value`, `isVerified`, `userId`.
    * Go: struct `BioMetrics` with field `OwnerId` instead of `UserId`; primary key named `BioMetrics_Id`.

5. **Order**

    * SQL: has `orderNumber`, `requirements`, `dueDate`, `completedAt`.
    * Go: has `Title` (no such column), no `requirements`, no `completedAt`, renamed `dueDate`→`Deadline`.

6. **Chat**

    * SQL: `Chat.id` PRIMARY KEY, columns `buyerId`, `sellerId`.
    * Go: `Chat` has composite PK on `ChatID` **and** `OrderId` (but SQL `Chat` has no `orderId`).

7. **Message**

    * SQL: single PK `id`; columns `chatId`, `senderId`, `isEdited`.
    * Go: composite PK on `MessageId`, `ChatId`, `SenderId`, `ReceiversId`; missing `isEdited`; adds `ReceiversId`.

8. **GigImage**

    * SQL: `GigImage` table.
    * Go: struct named `GiGImage` (typo in capitalization) and foreign key `RelatedGig` rather than `Gig`.

9. **Gig**

    * SQL: PK `id`; FKs `categoryId`, `sellerId`.
    * Go: marks `CategoryId` and `SellerId` as part of composite PK (wrong) instead of simple FKs.

10. **GigPackage**

    * SQL: PK `id`; FK `gigId`.
    * Go: makes `GigId` part of composite PK (wrong), and field `Revisions` as `float64` vs SQL integer.

11. **GigTag and \_GigToGigTag**

    * SQL: many‑to‑many via `_GigToGigTag(A, B)`.
    * Go: `GigToGigTag` struct uses fields `A` & `B`, but foreign key tags refer to nonexistent `TOTAGId`.

12. **Category**

    * SQL: also has `parentId` for self‑reference.
    * Go: no `ParentID` field in `Category` struct.

13. **Review**

    * SQL: PK `id`; unique on `orderId`; FKs `authorId`, `orderId`; boolean `isPublic`.
    * Go: composite PK on `ReviewId`, `AuthorId`, `OrderId`; no `isPublic`; changes `rating` from integer to decimal.

---

**Summary:**

* **3 tables** have **no corresponding** Go model: **Badge**, **UserBadge**, **RegistrationToken**.
* **1 Go model** has **no SQL table**: **Country**.
* **Multiple composite‑PK and field‑type mismatches** across almost every entity with foreign keys.
* **Self‑references** (e.g. `Category.parentId`) and **optional columns** (e.g. SQL `completedAt`) are unmodeled.

These discrepancies will prevent GORM from correctly mapping associations and may lead to schema drift or runtime errors.
