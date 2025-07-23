#!/usr/bin/env python3
"""
Test script for CurseForge Auto-Updater
"""

import sys
import os
from pathlib import Path

# Add the package to the path
sys.path.insert(0, str(Path(__file__).parent))

def test_imports():
    """Test that all modules can be imported."""
    print("🧪 Testing imports...")
    
    try:
        import updater
        print("✅ updater imported successfully")
        
        from updater import main, get_config, api, downloader, utils
        print("✅ All submodules imported successfully")
        
        # Test configuration
        config = get_config()
        print(f"✅ Configuration loaded (API key present: {bool(config['api_key'])})")
        
        return True
        
    except ImportError as e:
        print(f"❌ Import error: {e}")
        return False
    except Exception as e:
        print(f"❌ Unexpected error: {e}")
        return False

def test_api_structure():
    """Test API module structure."""
    print("\n🧪 Testing API module structure...")
    
    try:
        from updater import api
        
        # Check that required functions exist
        required_functions = [
            'get_mod_info', 'get_mod_files', 'validate_api_key',
            '_make_request', 'CurseForgeAPIError'
        ]
        
        for func_name in required_functions:
            if hasattr(api, func_name):
                print(f"✅ {func_name} exists")
            else:
                print(f"❌ {func_name} missing")
                return False
        
        return True
        
    except Exception as e:
        print(f"❌ Error testing API structure: {e}")
        return False

def test_environment():
    """Test environment setup."""
    print("\n🧪 Testing environment...")
    
    # Check for .env file
    env_file = Path(".env")
    env_example = Path(".env.example")
    
    if env_file.exists():
        print("✅ .env file exists")
    elif env_example.exists():
        print("⚠️  .env file missing, but .env.example exists")
        print("   Run: cp .env.example .env")
    else:
        print("❌ No .env or .env.example file found")
        return False
    
    # Check for required directories
    downloads_dir = Path("downloads")
    if not downloads_dir.exists():
        print("📁 Creating downloads directory...")
        downloads_dir.mkdir(exist_ok=True)
        print("✅ Downloads directory created")
    else:
        print("✅ Downloads directory exists")
    
    return True

def main():
    """Run all tests."""
    print("🚀 CurseForge Auto-Updater Test Suite")
    print("=" * 50)
    
    tests = [
        test_imports,
        test_api_structure, 
        test_environment
    ]
    
    passed = 0
    total = len(tests)
    
    for test in tests:
        if test():
            passed += 1
        print()
    
    print("📊 Test Results:")
    print(f"   Passed: {passed}/{total}")
    
    if passed == total:
        print("✅ All tests passed!")
        
        # Show next steps
        print("\n🎯 Next steps:")
        print("1. Copy .env.example to .env (if not done)")
        print("2. Add your CurseForge API key to .env")
        print("3. Run: python3 -m updater.main")
        print("   Or: python cli.py")
        
        return 0
    else:
        print("❌ Some tests failed!")
        return 1

if __name__ == "__main__":
    sys.exit(main())
