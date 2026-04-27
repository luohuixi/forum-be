# posts 表读取清单

下面列的是整个项目里**会从数据库读取 `posts` 表**的功能，按“直接读表 / 间接读表”区分。

## 1. 直接读取 `posts` 表的功能

- **获取单篇帖子**
  - 路径：`microservice/gateway/handler/post/get.go` → `PostClient.GetPost`
  - 服务：`microservice/post/service/getPost.go`
  - DAO：`microservice/post/dao/post.go:GetPostInfo`
  - 读法：`Table("posts") ... First/Select/Joins`

- **首页帖子列表**
  - 路径：`microservice/gateway/handler/post/listMainPost.go` → `PostClient.ListMainPost`
  - 服务：`microservice/post/service/listMainPost.go`
  - DAO：`microservice/post/dao/post.go:ListMainPost`
  - 读法：`Table("posts") ... Scan`，并 join `users`

- **用户发布的帖子列表**
  - 路径：`microservice/gateway/handler/post/listUserPost.go` → `PostClient.ListUserCreatedPost`
  - 服务：`microservice/post/service/listUserCreatedPost.go`
  - DAO：`microservice/post/dao/post.go:ListUserCreatedPost` + `ListPostInfoByPostIds`
  - 读法：先查 `posts.id`，再按 ID 批量读 `posts`

- **收藏列表**
  - 路径：`microservice/gateway/handler/collection/list.go` → `PostClient.ListCollection`
  - 服务：`microservice/post/service/listCollection.go`
  - DAO：`microservice/post/dao/collection.go:ListCollectionByUserId` + `post.go:ListPostInfoByPostIds`
  - 读法：先查 `collections`，再批量读 `posts`

- **点赞列表**
  - 路径：`microservice/gateway/handler/like/getUserLikeList.go` → `PostClient.ListLikeByUserId`
  - 服务：`microservice/post/service/listLikeByUserId.go`
  - DAO：`microservice/post/dao/like.go:ListUserLike` + `post.go:ListPostInfoByPostIds`
  - 读法：先查 Redis 点赞列表，再批量读 `posts`

- **未读帖子数**
  - 路径：`microservice/gateway/handler/post/getUnReadPostNum.go` → `PostClient.GetUnReadPostNum`
  - 服务：`microservice/post/service/getUnReadPostNum.go`
  - DAO：`microservice/post/dao/post.go:CountPostByTime`
  - 读法：按时间/分类统计 `posts`

- **热门/质量/分类列表（统一入口）**
  - 路径：`microservice/gateway/handler/post/listMainPost.go`
  - 服务：`microservice/post/service/listMainPost.go`
  - DAO：`microservice/post/dao/post.go:ListMainPost`
  - 读法：`posts` 表按条件筛选、排序、分页

## 2. 写接口里顺带读取 `posts` 表的功能

- **发评论**
  - 路径：`microservice/post/service/createComment.go`
  - 读法：`GetPost(req.PostId)` 先读 `posts`，再创建 comment

- **收藏/取消收藏帖子**
  - 路径：`microservice/post/service/createOrRemoveCollection.go`
  - 读法：先通过 `IsUserCollectionPost` 判断，再会回查帖子信息用于响应

- **提交举报**
  - 路径：`microservice/post/service/createReport.go`
  - 读法：`req.TypeName == post` 时先 `GetPost(req.Id)`

- **更新帖子信息**
  - 路径：`microservice/post/service/updatePostInfo.go`
  - 读法：先 `GetPost(req.Id)`，再更新 `posts`

- **删除帖子 / 删除收藏 / 设精华**
  - 路径：`microservice/post/service/deleteItem.go`
  - 读法：先删除业务数据，再调用 Casbin；删除帖子本身时会先 `DeletePost`，内部会读 `posts`

## 3. 间接读取 `posts` 表的 DAO

- **举报列表/举报处理**
  - 文件：`microservice/post/dao/report.go`
  - 读法：
    - `ValidReport` 会在处理 `post` 举报时调用 `DeletePost`
    - `InValidReport` 会在取消自动封禁时调用 `GetPost`
    - `ListReport` 对 `post` 举报会 join `posts` 和 `users`

- **标签相关**
  - 文件：`microservice/post/dao/tag.go`
  - 读法：
    - `ListTagsByPostId` 通过 `post2tags` 找帖子标签
    - `isExistPostWithTagIdAndCategory` 会 join `posts` 判断某标签在某分类是否还存在

- **帖子详情聚合**
  - 文件：`microservice/post/dao/post.go`
  - 读法：`ListPostInfoByPostIds` 会继续读取点赞、评论、收藏、标签，属于一层帖子批量聚合读取

## 4. 直接读 `posts` 的核心 DAO 方法

- `PostModel.Get`
- `GetPostInfo`
- `ListMainPost`
- `ListUserCreatedPost`
- `ListPostInfoByPostIds`
- `CountPostByTime`
- `syncPostScore` / `syncItemLike`（定时同步时也会更新 `posts`）

## 5. 结论

如果只看“帖子内容/列表/详情”，核心入口就是：

1. `GetPost`
2. `ListMainPost`
3. `ListUserCreatedPost`
4. `ListCollection`
5. `ListLikeByUserId`
6. `GetUnReadPostNum`

但如果把“写接口里先读帖子再处理”也算上，还要额外包含：

- `CreateComment`
- `CreateOrRemoveCollection`
- `CreateReport`
- `UpdatePostInfo`
