-- CreateTable
CREATE TABLE "Badge" (
                         "id" TEXT NOT NULL PRIMARY KEY,
                         "label" TEXT NOT NULL,
                         "icon" TEXT NOT NULL,
                         "color" TEXT NOT NULL,
                         "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         "updatedAt" DATETIME NOT NULL
);

-- CreateTable
CREATE TABLE "UserBadge" (
                             "id" TEXT NOT NULL PRIMARY KEY,
                             "userId" TEXT NOT NULL,
                             "badgeId" TEXT NOT NULL,
                             "tier" TEXT NOT NULL DEFAULT 'BRONZE',
                             "isFeatured" BOOLEAN NOT NULL DEFAULT false,
                             "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                             "updatedAt" DATETIME NOT NULL,
                             CONSTRAINT "UserBadge_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
                             CONSTRAINT "UserBadge_badgeId_fkey" FOREIGN KEY ("badgeId") REFERENCES "Badge" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- CreateTable
CREATE TABLE "User" (
                        "id" TEXT NOT NULL PRIMARY KEY,
                        "firstName" TEXT NOT NULL,
                        "lastName" TEXT NOT NULL,
                        "email" TEXT NOT NULL,
                        "password" TEXT NOT NULL,
                        "verified" BOOLEAN NOT NULL DEFAULT false,
                        "twoFactorVerified" BOOLEAN NOT NULL DEFAULT false,
                        "username" TEXT NOT NULL,
                        "avatar" TEXT,
                        "country" TEXT NOT NULL,
                        "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        "updatedAt" DATETIME NOT NULL,
                        "walletCreated" BOOLEAN NOT NULL DEFAULT false,
                        "walletCreatedTime" DATETIME
);

-- CreateTable
CREATE TABLE "Skill" (
                         "id" TEXT NOT NULL PRIMARY KEY,
                         "label" TEXT NOT NULL,
                         "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         "updatedAt" DATETIME NOT NULL
);

-- CreateTable
CREATE TABLE "GigTag" (
                          "id" TEXT NOT NULL PRIMARY KEY,
                          "label" TEXT NOT NULL,
                          "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          "updatedAt" DATETIME NOT NULL
);

-- CreateTable
CREATE TABLE "user_skills" (
                               "level" INTEGER NOT NULL DEFAULT 1,
                               "endorsed" BOOLEAN NOT NULL DEFAULT false,
                               "skillId" TEXT NOT NULL,
                               "userId" TEXT NOT NULL,
                               "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                               "updatedAt" DATETIME NOT NULL,

                               PRIMARY KEY ("skillId", "userId"),
                               CONSTRAINT "user_skills_skillId_fkey" FOREIGN KEY ("skillId") REFERENCES "Skill" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
                               CONSTRAINT "user_skills_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- CreateTable
CREATE TABLE "Biometrics" (
                              "id" TEXT NOT NULL PRIMARY KEY,
                              "type" TEXT NOT NULL,
                              "value" TEXT NOT NULL,
                              "isVerified" BOOLEAN NOT NULL DEFAULT false,
                              "userId" TEXT NOT NULL,
                              "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                              "updatedAt" DATETIME NOT NULL,
                              CONSTRAINT "Biometrics_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- CreateTable
CREATE TABLE "Gig" (
                       "id" TEXT NOT NULL PRIMARY KEY,
                       "title" TEXT NOT NULL,
                       "description" TEXT NOT NULL,
                       "isActive" BOOLEAN NOT NULL DEFAULT true,
                       "viewCount" INTEGER NOT NULL DEFAULT 0,
                       "averageRating" REAL NOT NULL DEFAULT 0,
                       "ratingCount" INTEGER NOT NULL DEFAULT 0,
                       "categoryId" TEXT NOT NULL,
                       "sellerId" TEXT NOT NULL,
                       "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       "updatedAt" DATETIME NOT NULL,
                       CONSTRAINT "Gig_sellerId_fkey" FOREIGN KEY ("sellerId") REFERENCES "User" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
                       CONSTRAINT "Gig_categoryId_fkey" FOREIGN KEY ("categoryId") REFERENCES "Category" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- CreateTable
CREATE TABLE "RegistrationToken" (
                                     "code" TEXT NOT NULL PRIMARY KEY,
                                     "email" TEXT NOT NULL,
                                     "expiresAt" DATETIME NOT NULL,
                                     "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     "updatedAt" DATETIME NOT NULL,
                                     "userId" TEXT,
                                     CONSTRAINT "RegistrationToken_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User" ("id") ON DELETE SET NULL ON UPDATE CASCADE
);

-- CreateTable
CREATE TABLE "GigImage" (
                            "id" TEXT NOT NULL PRIMARY KEY,
                            "url" TEXT NOT NULL,
                            "alt" TEXT,
                            "isPrimary" BOOLEAN NOT NULL DEFAULT false,
                            "sortOrder" INTEGER NOT NULL DEFAULT 0,
                            "gigId" TEXT NOT NULL,
                            "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            "updatedAt" DATETIME NOT NULL,
                            CONSTRAINT "GigImage_gigId_fkey" FOREIGN KEY ("gigId") REFERENCES "Gig" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- CreateTable
CREATE TABLE "GigPackage" (
                              "id" TEXT NOT NULL PRIMARY KEY,
                              "title" TEXT NOT NULL,
                              "description" TEXT NOT NULL,
                              "price" REAL NOT NULL,
                              "deliveryTime" INTEGER NOT NULL,
                              "revisions" INTEGER NOT NULL DEFAULT 1,
                              "isActive" BOOLEAN NOT NULL DEFAULT true,
                              "gigId" TEXT NOT NULL,
                              "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                              "updatedAt" DATETIME NOT NULL,
                              CONSTRAINT "GigPackage_gigId_fkey" FOREIGN KEY ("gigId") REFERENCES "Gig" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- CreateTable
CREATE TABLE "GigPackageFeature" (
                                     "id" TEXT NOT NULL PRIMARY KEY,
                                     "title" TEXT NOT NULL,
                                     "description" TEXT,
                                     "included" BOOLEAN NOT NULL DEFAULT true,
                                     "gigPackageId" TEXT NOT NULL,
                                     "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     "updatedAt" DATETIME NOT NULL,
                                     CONSTRAINT "GigPackageFeature_gigPackageId_fkey" FOREIGN KEY ("gigPackageId") REFERENCES "GigPackage" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- CreateTable
CREATE TABLE "Review" (
                          "id" TEXT NOT NULL PRIMARY KEY,
                          "title" TEXT NOT NULL,
                          "content" TEXT,
                          "rating" INTEGER NOT NULL DEFAULT 5,
                          "isPublic" BOOLEAN NOT NULL DEFAULT true,
                          "authorId" TEXT NOT NULL,
                          "orderId" TEXT NOT NULL,
                          "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          "updatedAt" DATETIME NOT NULL,
                          CONSTRAINT "Review_orderId_fkey" FOREIGN KEY ("orderId") REFERENCES "Order" ("id") ON DELETE RESTRICT ON UPDATE CASCADE,
                          CONSTRAINT "Review_authorId_fkey" FOREIGN KEY ("authorId") REFERENCES "User" ("id") ON DELETE RESTRICT ON UPDATE CASCADE
);

-- CreateTable
CREATE TABLE "Category" (
                            "id" TEXT NOT NULL PRIMARY KEY,
                            "label" TEXT NOT NULL,
                            "slug" TEXT NOT NULL,
                            "isActive" BOOLEAN NOT NULL DEFAULT true,
                            "sortOrder" INTEGER NOT NULL DEFAULT 0,
                            "parentId" TEXT,
                            "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            "updatedAt" DATETIME NOT NULL,
                            CONSTRAINT "Category_parentId_fkey" FOREIGN KEY ("parentId") REFERENCES "Category" ("id") ON DELETE SET NULL ON UPDATE CASCADE
);

-- CreateTable
CREATE TABLE "Order" (
                         "id" TEXT NOT NULL PRIMARY KEY,
                         "orderNumber" TEXT NOT NULL,
                         "price" REAL NOT NULL,
                         "paymentMethod" TEXT NOT NULL,
                         "status" TEXT NOT NULL DEFAULT 'PENDING',
                         "transactionId" TEXT,
                         "requirements" TEXT,
                         "packageId" TEXT NOT NULL,
                         "sellerId" TEXT NOT NULL,
                         "buyerId" TEXT NOT NULL,
                         "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         "updatedAt" DATETIME NOT NULL,
                         "completedAt" DATETIME,
                         "dueDate" DATETIME,
                         CONSTRAINT "Order_sellerId_fkey" FOREIGN KEY ("sellerId") REFERENCES "User" ("id") ON DELETE RESTRICT ON UPDATE CASCADE,
                         CONSTRAINT "Order_packageId_fkey" FOREIGN KEY ("packageId") REFERENCES "GigPackage" ("id") ON DELETE RESTRICT ON UPDATE CASCADE,
                         CONSTRAINT "Order_buyerId_fkey" FOREIGN KEY ("buyerId") REFERENCES "User" ("id") ON DELETE RESTRICT ON UPDATE CASCADE
);

-- CreateTable
CREATE TABLE "Chat" (
                        "id" TEXT NOT NULL PRIMARY KEY,
                        "lastActivity" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        "isArchived" BOOLEAN NOT NULL DEFAULT false,
                        "buyerId" TEXT NOT NULL,
                        "sellerId" TEXT NOT NULL,
                        "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        "updatedAt" DATETIME NOT NULL,
                        CONSTRAINT "Chat_sellerId_fkey" FOREIGN KEY ("sellerId") REFERENCES "User" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
                        CONSTRAINT "Chat_buyerId_fkey" FOREIGN KEY ("buyerId") REFERENCES "User" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- CreateTable
CREATE TABLE "Message" (
                           "id" TEXT NOT NULL PRIMARY KEY,
                           "content" TEXT NOT NULL,
                           "isRead" BOOLEAN NOT NULL DEFAULT false,
                           "isEdited" BOOLEAN NOT NULL DEFAULT false,
                           "senderId" TEXT NOT NULL,
                           "chatId" TEXT NOT NULL,
                           "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           "updatedAt" DATETIME NOT NULL,
                           CONSTRAINT "Message_chatId_fkey" FOREIGN KEY ("chatId") REFERENCES "Chat" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
                           CONSTRAINT "Message_senderId_fkey" FOREIGN KEY ("senderId") REFERENCES "User" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- CreateTable
CREATE TABLE "_GigToGigTag" (
                                "A" TEXT NOT NULL,
                                "B" TEXT NOT NULL,
                                CONSTRAINT "_GigToGigTag_A_fkey" FOREIGN KEY ("A") REFERENCES "Gig" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
                                CONSTRAINT "_GigToGigTag_B_fkey" FOREIGN KEY ("B") REFERENCES "GigTag" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- CreateIndex
CREATE UNIQUE INDEX "Badge_label_key" ON "Badge"("label");

-- CreateIndex
CREATE UNIQUE INDEX "UserBadge_userId_badgeId_key" ON "UserBadge"("userId", "badgeId");

-- CreateIndex
CREATE UNIQUE INDEX "User_email_key" ON "User"("email");

-- CreateIndex
CREATE UNIQUE INDEX "User_username_key" ON "User"("username");

-- CreateIndex
CREATE INDEX "User_email_idx" ON "User"("email");

-- CreateIndex
CREATE INDEX "User_username_idx" ON "User"("username");

-- CreateIndex
CREATE UNIQUE INDEX "Skill_label_key" ON "Skill"("label");

-- CreateIndex
CREATE UNIQUE INDEX "GigTag_label_key" ON "GigTag"("label");

-- CreateIndex
CREATE UNIQUE INDEX "user_skills_skillId_userId_key" ON "user_skills"("skillId", "userId");

-- CreateIndex
CREATE UNIQUE INDEX "Biometrics_userId_type_key" ON "Biometrics"("userId", "type");

-- CreateIndex
CREATE INDEX "Gig_sellerId_idx" ON "Gig"("sellerId");

-- CreateIndex
CREATE INDEX "Gig_isActive_idx" ON "Gig"("isActive");

-- CreateIndex
CREATE UNIQUE INDEX "RegistrationToken_email_key" ON "RegistrationToken"("email");

-- CreateIndex
CREATE INDEX "GigImage_gigId_isPrimary_idx" ON "GigImage"("gigId", "isPrimary");

-- CreateIndex
CREATE INDEX "GigPackage_gigId_isActive_idx" ON "GigPackage"("gigId", "isActive");

-- CreateIndex
CREATE UNIQUE INDEX "Review_orderId_key" ON "Review"("orderId");

-- CreateIndex
CREATE INDEX "Review_rating_idx" ON "Review"("rating");

-- CreateIndex
CREATE UNIQUE INDEX "Category_label_key" ON "Category"("label");

-- CreateIndex
CREATE UNIQUE INDEX "Category_slug_key" ON "Category"("slug");

-- CreateIndex
CREATE INDEX "Category_isActive_sortOrder_idx" ON "Category"("isActive", "sortOrder");

-- CreateIndex
CREATE UNIQUE INDEX "Order_orderNumber_key" ON "Order"("orderNumber");

-- CreateIndex
CREATE INDEX "Order_status_idx" ON "Order"("status");

-- CreateIndex
CREATE INDEX "Order_buyerId_status_idx" ON "Order"("buyerId", "status");

-- CreateIndex
CREATE INDEX "Order_sellerId_status_idx" ON "Order"("sellerId", "status");

-- CreateIndex
CREATE INDEX "Chat_lastActivity_idx" ON "Chat"("lastActivity");

-- CreateIndex
CREATE UNIQUE INDEX "Chat_sellerId_buyerId_key" ON "Chat"("sellerId", "buyerId");

-- CreateIndex
CREATE INDEX "Message_chatId_createdAt_idx" ON "Message"("chatId", "createdAt");

-- CreateIndex
CREATE INDEX "Message_isRead_idx" ON "Message"("isRead");

-- CreateIndex
CREATE UNIQUE INDEX "_GigToGigTag_AB_unique" ON "_GigToGigTag"("A", "B");

-- CreateIndex
CREATE INDEX "_GigToGigTag_B_index" ON "_GigToGigTag"("B");