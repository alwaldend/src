//! Generate command identifiers through Anki's own descriptor helpers.

use std::{env, fs, path::PathBuf};

/// main generates a command-name lookup using the same indices as rslib.
fn main() {
    let path = anki_proto_gen::descriptors_path();
    println!("cargo:rerun-if-changed={}", path.display());
    let descriptors = fs::read(path).expect("read upstream Anki descriptors");
    let pool = prost_reflect::DescriptorPool::decode(descriptors.as_slice())
        .expect("parse upstream Anki descriptors");
    let (collection_services, services) = anki_proto_gen::get_services(&pool);
    let mut source =
        String::from("fn command_ids(name: &str) -> Option<(u32, u32)> { match name {\n");
    for service in services {
        for method in service.all_methods() {
            source.push_str(&format!(
                "{:?} => Some(({}, {})),\n",
                format!("/{}/{}", service.proto.full_name(), method.proto.name()),
                service.index,
                method.index
            ));
        }
        let collection_service = collection_services
            .iter()
            .find(|candidate| candidate.name == service.name.trim_start_matches("Backend"))
            .expect("matching upstream collection service");
        for method in &collection_service.trait_methods {
            let backend_method = service
                .all_methods()
                .find(|candidate| candidate.name == method.name)
                .expect("collection method exposed by Anki's backend");
            assert_eq!(method.proto.input(), backend_method.proto.input());
            assert_eq!(method.proto.output(), backend_method.proto.output());
            source.push_str(&format!(
                "{:?} => Some(({}, {})),\n",
                format!(
                    "/{}/{}",
                    collection_service.proto.full_name(),
                    method.proto.name()
                ),
                service.index,
                backend_method.index
            ));
        }
    }
    source.push_str("_ => None, } }\n");
    fs::write(
        PathBuf::from(env::var("OUT_DIR").expect("Cargo output directory")).join("commands.rs"),
        source,
    )
    .expect("write generated command lookup");
}
