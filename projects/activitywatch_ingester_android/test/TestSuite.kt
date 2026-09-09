package com.alwaldend.src.projects.activitywatch_ingester_android

import org.junit.runner.RunWith
import org.junit.runners.Suite

@RunWith(Suite::class)
@Suite.SuiteClasses(
    ActivityWatchEventTest::class,
    ActivityWatchIngesterTest::class,
)
class TestSuite
