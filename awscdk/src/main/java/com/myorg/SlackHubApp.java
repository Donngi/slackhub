package com.myorg;

import software.amazon.awscdk.core.App;

public class SlackHubApp {
    public static void main(final String[] args) {
        App app = new App();
        
        new SlackHubStack(app, "SlackHubStack");

        app.synth();
    }
}
