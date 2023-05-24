git reset "$(git merge-base master "$(git branch --show-current)")"
git add -A && git commit -m 'v0.0.1'
git push --force
# git remote set-url origin https://<你的令牌>@github.com/<你的git用户名>/<要修改的仓库名>.git