```json lines
{
  "Users": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@example.com",
      "password": "$2b$10$...",
      "verified": true,
      "twoFactorVerified": false,
      "username": "johndoe123",
      "avatar": "https://cdn.example.com/avatar1.jpg",
      "country": "US",
      "walletCreated": true,
      "walletCreatedTime": "2023-01-15T08:00:00Z",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-15T08:30:00Z",
      "badges": [
        {
          "id": "badge1-uuid",
          "label": "Top Seller",
          "icon": "star.svg",
          "color": "gold",
          "tier": "GOLD",
          "isFeatured": true,
          "createdAt": "2023-01-10T00:00:00Z"
        }
      ],
      "skills": [
        {
          "skillId": "skill1-uuid",
          "label": "Web Development",
          "level": 3,
          "endorsed": true
        }
      ],
      "biometrics": [
        {
          "id": "bio1-uuid",
          "type": "FACE_ID",
          "value": "face_encrypted_data",
          "isVerified": true
        }
      ]
    }
  ],
  
  "Gigs": [
    {
      "id": "gig1-uuid",
      "title": "Professional Website Development",
      "description": "Full-stack web development services",
      "isActive": true,
      "viewCount": 1500,
      "averageRating": 4.8,
      "ratingCount": 45,
      "categoryId": "cat1-uuid",
      "sellerId": "550e8400-e29b-41d4-a716-446655440000",
      "packages": [
        {
          "id": "pkg1-uuid",
          "title": "Standard Package",
          "description": "5-page website with CMS",
          "price": 500.00,
          "deliveryTime": 14,
          "revisions": 3,
          "features": [
            {
              "id": "feat1-uuid",
              "title": "Responsive Design",
              "included": true
            }
          ]
        }
      ],
      "images": [
        {
          "url": "https://cdn.example.com/gig1.jpg",
          "alt": "Website Example",
          "isPrimary": true
        }
      ],
      "tags": ["web", "development"]
    }
  ],

  "Orders": [
    {
      "id": "order1-uuid",
      "order_number": "ORD-2023-001",
      "price": 500.00,
      "payment_method": "CREDIT_CARD",
      "status": "COMPLETED",
      "transaction_id": "txn_123456",
      "requirements": "Include contact form",
      "package_id": "pkg1-uuid",
      "seller_id": "550e8400-e29b-41d4-a716-446655440000",
      "buyer_id": "buyer1-uuid",
      "created_at": "2023-02-01T09:00:00Z",
      "completed_at": "2023-02-15T09:00:00Z"
    }
  ],

  "Reviews": [
    {
      "id": "rev1-uuid",
      "title": "Excellent Service!",
      "rating": 5,
      "content": "Delivered exactly what was promised",
      "authorId": "buyer1-uuid",
      "orderId": "order1-uuid"
    }
  ],

  "Categories": [
    {
      "id": "cat1-uuid",
      "label": "Web Development",
      "slug": "web-development",
      "isActive": true,
      "sortOrder": 1
    }
  ],

  "Badges": [
    {
      "id": "badge1-uuid",
      "label": "Top Seller",
      "icon": "star.svg",
      "color": "gold"
    }
  ],

  "Skills": [
    {
      "id": "skill1-uuid",
      "label": "Web Development"
    }
  ]
}
```