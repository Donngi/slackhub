package com.github.nicojn.slackhub.sample.sample;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpRequest.BodyPublishers;
import java.net.http.HttpResponse.BodyHandlers;

import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestStreamHandler;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

public class SampleToolHandler implements RequestStreamHandler {

	@Override
	public void handleRequest(InputStream input, OutputStream output, Context context) throws IOException {
        
        // Read input json value
		var byteStream = new ByteArrayOutputStream();
		byte[] buf = new byte[1024];
		int length;
		while ((length = input.read(buf)) != -1) {
		    byteStream.write(buf, 0, length);
		}
		var strInput = byteStream.toString("UTF-8");

        var mapper = new ObjectMapper();
        JsonNode node = mapper.readTree(strInput);
		
		// Get user input values.
		var lunchInput = node.get("view").get("state").get("values").get("lunch_block").get("lunch_action").get("value").asText();
        var detailInput = node.get("view").get("state").get("values").get("detail_block").get("detail_action").get("value").asText();
        
        // Get private_metadata fields value.
        // NOTE: private_metadata field's type is simply "String". SlackHub stores some convenient values in the form of JSON.
        // In order to use them easily, you should convert private_metadata field value to JSON.
        var privateMetaStr = node.get("view").get("private_metadata").asText();
        JsonNode privateMetaNode = mapper.readTree(privateMetaStr);
        var responseUrl = privateMetaNode.get("response_url").asText();
        
        // Send a message to Slack
        var msg = "What is your favorite lunch?\n - " + lunchInput + "\nTell us more!\n - " + detailInput;
        var param = "{\"text\": \"" + msg + "\"}";
        
        var request = HttpRequest
				.newBuilder(URI.create(responseUrl))
				.POST(BodyPublishers.ofString(param))
				.setHeader("Content-Type", "application/json")
				.build();
        
        var httpClient = HttpClient.newBuilder()
        		.build();
        try {
			httpClient.send(request, BodyHandlers.ofString());
		} catch (IOException | InterruptedException e) {
			e.printStackTrace();
		}
        
        return;
    }

}