package com.sourcegraph.telemetry;

import com.google.gson.JsonObject;
import org.jetbrains.annotations.NotNull;
import org.jetbrains.annotations.Nullable;

public class Event {
    String eventName;
    String anonymousUserId;
    String url;
    JsonObject eventProperties;
    /**
     * PRIVACY: Do NOT include any potentially private information, such as search queries or repository names.
     */
    JsonObject publicArgument;
    int eventId;

    public Event(@NotNull String eventName,
                 @NotNull String anonymousUserId,
                 @NotNull String url,
                 @Nullable JsonObject eventProperties,
                 @Nullable JsonObject publicArgument) {
        this.eventName = eventName;
        this.anonymousUserId = anonymousUserId;
        this.url = url;
        this.eventProperties = eventProperties;
        this.publicArgument = publicArgument;
    }

    public JsonObject toJson() {
        JsonObject returnValue = new JsonObject();
        returnValue.addProperty("event", this.eventName);
        returnValue.addProperty("userCookieID", this.anonymousUserId);
        returnValue.addProperty("url", this.url);
        returnValue.addProperty("source", "IDEEXTENSION");
        returnValue.addProperty("referrer", "JETBRAINS");
        if (eventProperties != null) {
            returnValue.add("argument", eventProperties);
        }
        if (publicArgument != null) {
            returnValue.add("publicArgument", publicArgument);
        }
        returnValue.addProperty("deviceID", this.anonymousUserId);
        return returnValue;
    }
}
