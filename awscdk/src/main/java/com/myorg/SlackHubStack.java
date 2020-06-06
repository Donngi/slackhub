package com.myorg;

import java.io.FileInputStream;
import java.io.IOException;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;
import java.util.Properties;

import software.amazon.awscdk.core.Construct;
import software.amazon.awscdk.core.Stack;
import software.amazon.awscdk.core.StackProps;
import software.amazon.awscdk.services.apigateway.LambdaRestApi;
import software.amazon.awscdk.services.dynamodb.Attribute;
import software.amazon.awscdk.services.dynamodb.AttributeType;
import software.amazon.awscdk.services.dynamodb.Table;
import software.amazon.awscdk.services.iam.Effect;
import software.amazon.awscdk.services.iam.IManagedPolicy;
import software.amazon.awscdk.services.iam.ManagedPolicy;
import software.amazon.awscdk.services.iam.PolicyStatement;
import software.amazon.awscdk.services.iam.Role;
import software.amazon.awscdk.services.iam.ServicePrincipal;
import software.amazon.awscdk.services.lambda.Code;
import software.amazon.awscdk.services.lambda.Function;
import software.amazon.awscdk.services.lambda.Runtime;

public class SlackHubStack extends Stack {
    public SlackHubStack(final Construct scope, final String id) {
        this(scope, id, null);
    }

    public SlackHubStack(final Construct scope, final String id, final StackProps props) {
        super(scope, id, props);

        // Load properties
        Properties properties = new Properties();
        String confPass = "../config.properties";

        try {
            FileInputStream stream = new FileInputStream(confPass);
            properties.load(stream);
        } catch (IOException e) {
            e.printStackTrace();
        }

        String region = properties.getProperty("region");
        String botKey = properties.getProperty("ssm-bot-user");
        String secretKey = properties.getProperty("ssm-signing-secret");

        // DynamoDB
        final Table tools = Table.Builder.create(this, "SlackHubToolsTable")
            .partitionKey(Attribute.builder()
                .name("id")
                .type(AttributeType.STRING)
                .build())
            .build();

        // Lamnda - sample (Go)
        final Function sampleGo = Function.Builder.create(this, "SlackHubSampleToolGo")
            .runtime(Runtime.GO_1_X)
            .code(Code.fromAsset("../examples/go/bin"))
            .handler("main")
            .build();

        // IAM Role for register
        IManagedPolicy ssmPolicy = ManagedPolicy.fromAwsManagedPolicyName("AmazonSSMReadOnlyAccess");
        Role registerRole = Role.Builder.create(this, "SlackHubRegisterRole")
            .managedPolicies(Arrays.asList(ssmPolicy))
            .assumedBy(new ServicePrincipal("lambda.amazonaws.com"))
            .build();

        registerRole.addToPolicy(PolicyStatement.Builder.create()
            .effect(Effect.ALLOW)
            .resources(Arrays.asList("*"))
            .actions(Arrays.asList("logs:CreateLogGroup","logs:CreateLogStream","logs:PutLogEvents"))
            .build());

        // Lamnda - register
        final Map<String,String> registerEnv = new HashMap<>();
        registerEnv.put("REGION", region);
        registerEnv.put("PARAMKEY_BOT_USER_AUTH_TOKEN", botKey);
        registerEnv.put("PARAMKEY_SIGNING_SECRET", secretKey);
        registerEnv.put("DYNAMO_DB_NAME", tools.getTableName());

        final Function register = Function.Builder.create(this, "SlackHubRegister")
            .runtime(Runtime.GO_1_X)
            .code(Code.fromAsset("../register/bin"))
            .handler("main")
            .environment(registerEnv)
            .role(registerRole)
            .build();

        tools.grantReadWriteData(register);

        // IAM Role for editor
        Role editorRole = Role.Builder.create(this, "SlackHubEditorRole")
            .managedPolicies(Arrays.asList(ssmPolicy))
            .assumedBy(new ServicePrincipal("lambda.amazonaws.com"))
            .build();
            
        editorRole.addToPolicy(PolicyStatement.Builder.create()
            .effect(Effect.ALLOW)
            .resources(Arrays.asList("*"))
            .actions(Arrays.asList("logs:CreateLogGroup","logs:CreateLogStream","logs:PutLogEvents"))
            .build());

        // Lamnda - editor
        final Map<String,String> editorEnv = new HashMap<>();
        editorEnv.put("REGION", region);
        editorEnv.put("PARAMKEY_BOT_USER_AUTH_TOKEN", botKey);
        editorEnv.put("PARAMKEY_SIGNING_SECRET", secretKey);
        editorEnv.put("DYNAMO_DB_NAME", tools.getTableName());

        final Function editor = Function.Builder.create(this, "SlackHubEditor")
            .runtime(Runtime.GO_1_X)
            .code(Code.fromAsset("../editor/bin"))
            .handler("main")
            .environment(editorEnv)
            .role(editorRole)
            .build();

        tools.grantReadWriteData(editor);

        // IAM Role for catalog
        Role catalogRole = Role.Builder.create(this, "SlackHubCatalogRole")
            .managedPolicies(Arrays.asList(ssmPolicy))
            .assumedBy(new ServicePrincipal("lambda.amazonaws.com"))
            .build();

        catalogRole.addToPolicy(PolicyStatement.Builder.create()
            .effect(Effect.ALLOW)
            .resources(Arrays.asList("*"))
            .actions(Arrays.asList("logs:CreateLogGroup","logs:CreateLogStream","logs:PutLogEvents"))
            .build());

        // Lamnda - catalog
        final Map<String,String> catalogEnv = new HashMap<>();
        catalogEnv.put("REGION", region);
        catalogEnv.put("PARAMKEY_BOT_USER_AUTH_TOKEN", botKey);
        catalogEnv.put("PARAMKEY_SIGNING_SECRET", secretKey);
        catalogEnv.put("DYNAMO_DB_NAME", tools.getTableName());

        final Function catalog = Function.Builder.create(this, "SlackHubCatalog")
            .runtime(Runtime.GO_1_X)
            .code(Code.fromAsset("../catalog/bin"))
            .handler("main")
            .environment(catalogEnv)
            .role(catalogRole)
            .build();

        tools.grantReadWriteData(catalog);

        // IAM Role for eraser
        Role eraserRole = Role.Builder.create(this, "SlackHubEraserRole")
            .managedPolicies(Arrays.asList(ssmPolicy))
            .assumedBy(new ServicePrincipal("lambda.amazonaws.com"))
            .build();

        eraserRole.addToPolicy(PolicyStatement.Builder.create()
            .effect(Effect.ALLOW)
            .resources(Arrays.asList("*"))
            .actions(Arrays.asList("logs:CreateLogGroup","logs:CreateLogStream","logs:PutLogEvents"))
            .build());

        // Lamnda - eraser
        final Map<String,String> eraserEnv = new HashMap<>();
        eraserEnv.put("REGION", region);
        eraserEnv.put("PARAMKEY_BOT_USER_AUTH_TOKEN", botKey);
        eraserEnv.put("PARAMKEY_SIGNING_SECRET", secretKey);
        eraserEnv.put("DYNAMO_DB_NAME", tools.getTableName());

        final Function eraser = Function.Builder.create(this, "SlackHubEraser")
            .runtime(Runtime.GO_1_X)
            .code(Code.fromAsset("../eraser/bin"))
            .handler("main")
            .environment(eraserEnv)
            .role(eraserRole)
            .build();

        tools.grantReadWriteData(eraser);

        // IAM Role for initial
        Role initialRole = Role.Builder.create(this, "SlackHubInitialRole")
            .managedPolicies(Arrays.asList(ssmPolicy))
            .assumedBy(new ServicePrincipal("lambda.amazonaws.com"))
            .build();
            
        initialRole.addToPolicy(PolicyStatement.Builder.create()
            .effect(Effect.ALLOW)
            .resources(Arrays.asList("*"))
            .actions(Arrays.asList("logs:CreateLogGroup","logs:CreateLogStream","logs:PutLogEvents"))
            .build());

        // Lambda - initial
        final Map<String,String> initialEnv = new HashMap<>();
        initialEnv.put("PARAMKEY_BOT_USER_AUTH_TOKEN", botKey);
        initialEnv.put("PARAMKEY_SIGNING_SECRET", secretKey);
        initialEnv.put("REGION", region);
        initialEnv.put("DYNAMO_DB_NAME", tools.getTableName());
        initialEnv.put("REGISTER_TOOL_ARN", register.getFunctionArn());
        initialEnv.put("EDITOR_TOOL_ARN", editor.getFunctionArn());
        initialEnv.put("CATALOG_TOOL_ARN", catalog.getFunctionArn());
        initialEnv.put("ERASER_TOOL_ARN", eraser.getFunctionArn());
        initialEnv.put("SAMPLE_TOOL_GO_ARN", sampleGo.getFunctionArn());

        final Function initial = Function.Builder.create(this, "SlackHubInitial")
            .runtime(Runtime.GO_1_X)
            .code(Code.fromAsset("../initial/bin"))
            .handler("main")
            .environment(initialEnv)
            .role(initialRole)
            .build();

        tools.grantReadWriteData(initial);

        // IAM Role for interactive
        Role interactiveRole = Role.Builder.create(this, "SlackHubInteractiveRole")
            .managedPolicies(Arrays.asList(ssmPolicy))
            .assumedBy(new ServicePrincipal("lambda.amazonaws.com"))
            .build();
            
        interactiveRole.addToPolicy(PolicyStatement.Builder.create()
            .effect(Effect.ALLOW)
            .resources(Arrays.asList("*"))
            .actions(Arrays.asList("logs:CreateLogGroup","logs:CreateLogStream","logs:PutLogEvents","lambda:InvokeFunction"))
            .build());

        // Lambda - interactive
        final Map<String,String> interactiveEnv = new HashMap<>();
        interactiveEnv.put("PARAMKEY_BOT_USER_AUTH_TOKEN", botKey);
        interactiveEnv.put("PARAMKEY_SIGNING_SECRET", secretKey);
        interactiveEnv.put("REGION", region);
        interactiveEnv.put("DYNAMO_DB_NAME", tools.getTableName());

        final Function interactive = Function.Builder.create(this, "SlackHubInteractive")
            .runtime(Runtime.GO_1_X)
            .code(Code.fromAsset("../interactive/bin"))
            .handler("main")
            .environment(interactiveEnv)
            .role(interactiveRole)
            .build();

        tools.grantReadWriteData(interactive);

        // API Gateway
        LambdaRestApi.Builder.create(this, "SlackHubInitialEndpoint")
            .handler(initial)
            .build();

        LambdaRestApi.Builder.create(this, "SlackHubInteractiveEndpoint")
            .handler(interactive)
            .build();
    }
}
