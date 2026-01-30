import unittest
from lib.runner import CommandRunner

class TestCommandRunner(unittest.TestCase):
    def test_run_step_dry_run(self):
        runner = CommandRunner(dry_run=True)
        success, output = runner.run_step({'cmd': 'echo', 'tag': 'test', 'friendly': 'Test'})
        self.assertTrue(success)
        self.assertEqual(output, "Dry run")

    def test_run_step_actual_run(self):
        runner = CommandRunner(dry_run=False)
        success, output = runner.run_step({'cmd': 'echo hello', 'tag': 'test', 'friendly': 'Test'})
        self.assertTrue(success)
        self.assertIn("hello", output)
