## Duplicate Freeze Window Error

When running the `DOE Federation CLI Utility` pipeline to copy an existing Harness project into a new organization, there may be time when you must run the pipeline more the one time, for example, when testing.  When this happens, you can run into an issue where the source project is already frozen and will result in an error on the second run. 

You'll notice an error when this situation arises, as shown below:

```bash
.... Copy error from pipeline
```

When this occurs because the `DOE Federation CLI Utility` is attempting to freeze the source project, but it has a freeze in place with the same name. To resolve this issue, you must remove the freeze from the source project.  You can do this by following the steps below:

1. Navigate to the source project in Harness.
2. Click on `Project Settings` in the left-hand menu.
3. Click on `Freeze Windows` under the "Security and Governance" section.
4. Find the freeze window named "hrns_copy_prj_freeze" and select the three dots on the right hand side.  Click on `Delete` to remove the freeze window.
5. Re-run the `DOE Federation CLI Utility` pipeline.

