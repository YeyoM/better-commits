# Better Commits

This is my version of a CLI tool that is meant to help you write better commits, this tool helps you write your commits based on the convetional commits guideline.

### How it works?

There are a total of 6 steps to create a conventional commit 

1. Select the type of change that you're committing 
2. What is the scope of this change (user, admin, etc.) the user can select none
3. Select a gitmoji for the commit, or the user can select none
4. Write a short, imperative tense description of the change
5. Provide a longer description of the change
6. List any breaking changes or issues closed by this change

### How to use it?  

Since this is a golang project, you can install it by cloning the repository and running the following command:

```bash
go build
```

This will create an executable file that you can run with the following command (Linux):

```bash
./better-commits
```

*Note, give the executable file the right permissions to run it*

After having the project compiled, you can run the following commands to add it to your path:

```bash
mv better-commits /usr/local/bin
```

Now you can run the command `better-commits` from anywhere in your terminal.

### How to contribute?

If you want to contribute to this project, you can fork the repository and create a pull request with your changes.

**Note:** This is my first golang project and is still in development, please be patient if there are some features that are not implemented yet, you can open an issue and together work on the feature or bug.

### License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.



