package core

import (
	"context"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

const StoreInputExample = `[
  {
    "Id": 507,
    "Content": "今天下雨本来心情差，结果点了一份老乡鸡卤肉拌饭，吃了心情好多了，狠狠推荐",
    "Title": "美食推荐"
  },
  {
    "Id": 508,
    "Content": "## **克隆仓库**
git clone
## **提交**
1. git add 
> git add <文件名>：将指定的文件添加到暂存区。
git add .：将所有修改过的文件添加到暂存区。
git add -A：将所有修改过的文件和新文件（包括未跟踪的文件）添加到暂存区
2. git commit -m
> ***git commit*** ：这将打开文本编辑器，让你输入提交信息。完成信息编写后保存并关闭编辑器，提交就会完成。
***git commit -m*** ：这是一种快速提交的方式，允许你直接在命令行中提供提交信息。例如，git commit -m "修复了登录功能的 bug" 会创建一个提交，其信息是“修复了登录功能的 bug”。
git commit 只影响本地仓库，并不会更改远程仓库（如 GitHub 上的仓库）。要将这些更改推送到远程仓库，你需要使用 git push 命令。
3. git push
> ***git push <远程仓库名> <分支名>*** ：这个命令会将指定的本地分支推送到指定的远程仓库。例如，git push origin master 会将本地的 master 分支推送到名为 origin 的远程仓库。
***git push*** ：如果已经设置了本地分支和远程分支之间的跟踪关系，可以直接使用这个命令来推送更改。Git 会自动推送到之前配置的远程分支。
***git push -u <远程仓库名> <分支名>*** ：除了推送更改外，这个命令还会设置本地分支和远程分支之间的跟踪关系。在首次推送分支时常用这个命令。例如，git push -u origin feature 会将本地的 feature 分支推送到远程仓库，并设置跟踪关系。",
    "Title": "git 基本操作"
  }
]`

const StoreOutputExample = `{
  "items": [
    {
      "content": "老乡鸡卤肉拌饭是一种让人心情变好的美食推荐。",
      "meta": {
        "post_id": "507",
        "title": "美食推荐"
      }
    },
    {
      "content": "Git 的三个基本命令分别是 git add、git commit 和 git push。git add 用于将修改加入暂存区，git commit 用于提交到本地仓库，git push 用于将本地提交推送到远程仓库。git push -u 用于建立本地分支和远程分支的跟踪关系。",
      "meta": {
        "post_id": "508",
        "title": "git 基本操作"
      }
    }
  ]
}`

func StorePrompt(content string) ([]*schema.Message, error) {
	input := prompt.FromMessages(
		schema.GoTemplate,
		schema.SystemMessage("你是一个论坛问答助手，你当前的任务是提取各帖子中有价值的内容，并将其存储到向量数据库中，存储有具体的工具可以调用，你需按照工具规范来"),
		schema.SystemMessage("示例输入: "+StoreInputExample),
		schema.SystemMessage("示例输出: "+StoreOutputExample),
		schema.UserMessage(content),
	)

	return input.Format(context.Background(), nil)
}
