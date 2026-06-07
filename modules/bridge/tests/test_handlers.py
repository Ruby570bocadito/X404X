"""X404X Python Bridge Handler Tests — 85 RPCs across 7 versions"""
import sys, os, json, unittest
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'handlers'))

class TestRansomwareBase(unittest.TestCase):
    def setUp(self):
        self.registry = {}
        try:
            import ransomware
            ransomware.register_routes(self.registry)
        except: pass
        try:
            import ransomware_advanced
            ransomware_advanced.register_routes(self.registry)
        except: pass
        try:
            import ransomware_blockz
            ransomware_blockz.register_routes(self.registry)
        except: pass
        try:
            import ransomware_v26
            ransomware_v26.register_routes(self.registry)
        except: pass
        try:
            import ransomware_v27
            ransomware_v27.register_routes(self.registry)
        except: pass
        try:
            import ransomware_v28
            ransomware_v28.register_routes(self.registry)
        except: pass
        try:
            import ransomware_v29
            ransomware_v29.register_routes(self.registry)
        except: pass
        try:
            import ransomware_v210
            ransomware_v210.register_routes(self.registry)
        except: pass

    def test_handler_count(self):
        count = sum(len(v) for v in self.registry.values())
        self.assertGreater(count, 50, f"Expected >50 handlers, got {count}")

    def test_all_handlers_are_callable(self):
        for group, handlers in self.registry.items():
            for name, handler in handlers.items():
                self.assertTrue(callable(handler), f"{group}.{name} not callable")

    def test_handlers_return_dict(self):
        for group, handlers in self.registry.items():
            for name, handler in handlers.items():
                result = handler({})
                self.assertIsInstance(result, dict, f"{group}.{name} did not return dict")

    def test_handlers_have_success_field(self):
        for group, handlers in self.registry.items():
            for name, handler in handlers.items():
                result = handler({})
                self.assertIn('success', result, f"{group}.{name} missing 'success'")

    def test_v210_apocalipsis(self):
        if 'ransomware_v210' in self.registry and 'apocalipsis' in self.registry['ransomware_v210']:
            r = self.registry['ransomware_v210']['apocalipsis']({})
            self.assertTrue(r['success'])

    def test_v210_phantom_evasion(self):
        if 'ransomware_v210' in self.registry and 'phantom_evasion' in self.registry['ransomware_v210']:
            r = self.registry['ransomware_v210']['phantom_evasion']({})
            self.assertTrue(r['success'])

if __name__ == '__main__':
    unittest.main(verbosity=2)
